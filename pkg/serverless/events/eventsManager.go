package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/config"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/workflow"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types/mqObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type EventsManager struct {
	route_table_manager            *RouteTableManager
	mqConn                         *minik8s_message.MQConnection
	func_request_frequency_manager *FuncRequestFrequencyManager
}

func NewEventsManager(mqConfig *minik8s_message.MQConfig) *EventsManager {
	fmt.Printf("New EventsManager\n")
	newConn, err := minik8s_message.NewMQConnection(mqConfig)
	if err != nil {
		fmt.Printf("New MQConnection Error, msg %s\n", err.Error())
		return nil
	}

	return &EventsManager{
		route_table_manager:            NewRouteTableManager(),
		mqConn:                         newConn,
		func_request_frequency_manager: NewFuncRequestFrequencyManager(),
	}
}

func (em *EventsManager) Init() {
	fmt.Printf("Init EventsManager\n")
	em.route_table_manager.Init()
	em.func_request_frequency_manager.Init()
	//注册函数指针
	config.GetFuncionPodRequestFrequency=em.getFuncionPodRequestFrequency
	config.TriggerServerlessFunction=em.triggerServerlessFunction
	config.TriggerServerlessWorkflow=em.triggerServerlessWorkflow
}

// func (em *EventsManager) Run() {
// 	fmt.Printf("Run EventsManager\n")
// }

// GetFuncionPodRequestFrequency 返回每个serverless function的每分钟每pod请求数，以便serving计算pod增减策略
// 如果一个serverless的pod数为0,则对其的扩容应该由events manager负责，不必上交给serving
//
//	@receiver em
//	@return map
func (em *EventsManager) getFuncionPodRequestFrequency() map[string]float64 {
	result := em.func_request_frequency_manager.GetAllRecentRequestFrequency()
	for funcname, funcfreq := range result {
		podnum := em.route_table_manager.GetFunctionPodNum(funcname)
		if podnum <= 0 {
			delete(result, funcname)
		} else {
			result[funcname] = funcfreq / float64(podnum)
		}
	}
	return result
}

func (em *EventsManager) SyncRouteTableRoutine() {
	// Never Return
	poller.PollerStaticPeriod(5*time.Second, em.route_table_manager.SyncRoutine, true)
}

// TriggerServerlessFunction 会触发一个函数的执行，在遇到0可用pod时请求创建，并向mq发送消息
//
//	@receiver em
//	@param funcName
//	@param params
//	@param mqName 	**可以为空字符串**，此时不会发送消息
//	@return string 	函数的返回值的序列化字符串
//	@return error
func (em *EventsManager) triggerServerlessFunction(funcName string, params string, mqName string) (string, error) {
	// 先为该函数加一条请求记录
	em.func_request_frequency_manager.AddOneRequest(funcName)

	// 在执行这个function时，外部需要确保此serverless函数已经存在
	/* faileCount记录失败次数，至多允许三次失败，如果发生三次失败，则认为请求失败 */
	failCount := 0
	/* 创建pod的轮数，初始开始时为-1，如果检查路由表为空，需要发送一次创建pod的请求
	请求创建后，每三轮检查路由表为空，则需要额外请求一次pod create，且加一次faileCount */
	podCreateCount := -1

	for {
		// 如果发生三次失败，则认为请求失败
		if failCount >= 3 {
			return "", errors.New("failCount reach three")
		}

		podIP := em.route_table_manager.FunctionName2PodIP(funcName)
		if podIP == ROUTETABLE_NONEPOD {
			if podCreateCount >= 3 {
				// 如果路由表中没有可用podIP,需要请求启动pod
				failCount += 1
				if failCount < 3 {
					em.requestNewFuncPod(funcName)
					podCreateCount = 0

					if mqName != "" {
						em.publishMessage(mqName, mqObject.MQmessage_Workflow{
							Isdone:        false,
							DataOrMessage: fmt.Sprintf("[Start Pod] There is no available pod for function %s, try to create...", funcName),
						})
					}
				}

			} else if podCreateCount == -1 {
				if mqName != "" {
					em.publishMessage(mqName, mqObject.MQmessage_Workflow{
						Isdone:        false,
						DataOrMessage: fmt.Sprintf("[Start Pod] There is no available pod for function %s, try to create...", funcName),
					})
				}
				em.requestNewFuncPod(funcName)
				podCreateCount = 0
			} else {
				podCreateCount += 1
			}
		} else {
			// 创建http的request/response数据结构和url,并请求
			requestbody := httpobject.HTTPRequest_callfunc{
				Params: params,
			}
			var responsebody httpobject.HTTPResponse_callfunc = httpobject.HTTPResponse_callfunc{}
			url := fmt.Sprintf(config.HTTPURL_callfunc_Template, podIP)
			status, err := httpRequest.PostRequestByObject(url, requestbody, &responsebody)
			if status != http.StatusOK || err != nil {
				// 请求错误
				failCount += 1
				fmt.Printf("routine error Post, status %d, return\n", status)
			} else {
				// 请求成功，返回数据
				return responsebody.Data, nil
			}
		}
		// 睡眠5s,等待各项数据的更新
		time.Sleep(5 * time.Second)
	}
}

func (em *EventsManager) triggerServerlessWorkflow(workflowObject workflow.Workflow, mqName string) {
	// 在执行这个workflow时，外部需要确保此workflow的所有serverless函数已经存在
	// FIXME: 如果workflow不能正常退出，引发线程僵死无法回收，如何解决？
	go em.workflowComputeRoutine(workflowObject, mqName)
}

func (em *EventsManager) requestNewFuncPod(funcName string) {
	requestbody := httpobject.HTTPRequest_AddServerlessFuncPod{
		FuncName: funcName,
	}
	status, err := httpRequest.PostRequestByObject(config.HTTPURL_AddServerlessFuncPod, requestbody, nil)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error Post, status %d, return\n", status)
		return
	}
}

func (em *EventsManager) workflowComputeRoutine(workflowObject workflow.Workflow, mqName string) {
	currentNode := workflowObject.Spec.WorkflowNodes[workflowObject.Spec.EntryNodeName]
	currentNodeName := workflowObject.Spec.EntryNodeName
	currentParams := workflowObject.Spec.Params
	step := 0
	var err error

	em.publishMessage(mqName, mqObject.MQmessage_Workflow{
		Isdone:        false,
		DataOrMessage: fmt.Sprintf("[START] The workflow %s is Triggered...", workflowObject.Metadata.Name),
	})

	is_end := false
	for {
		switch currentNode.NodeType {
		case workflow.SERVERLESS_NODETYPE_BRANCH:
			{
				step += 1
				em.publishMessage(mqName, mqObject.MQmessage_Workflow{
					Isdone:        false,
					DataOrMessage: fmt.Sprintf("[BRANCH] Step %d, node name %s...", step, currentNodeName),
				})
				branchResult, err := em.triggerServerlessFunction(currentNode.FunctionName, currentParams, mqName)
				if err != nil {
					em.publishMessage(mqName, mqObject.MQmessage_Workflow{
						Isdone:        true,
						DataOrMessage: fmt.Sprintf("[FAIL] Step %d, node name %s...", step, currentNodeName),
					})
					return
				}
				num, err := strconv.Atoi(branchResult)
				if err != nil {
					fmt.Printf("can't branch by Node Return, msg %s\n", err.Error())
					return
				}
				currentNodeName = currentNode.Branch[num]
				currentNode = workflowObject.Spec.WorkflowNodes[currentNodeName]
				break
			}
		case workflow.SERVERLESS_NODETYPE_CALCULATION:
			{
				step += 1
				em.publishMessage(mqName, mqObject.MQmessage_Workflow{
					Isdone:        false,
					DataOrMessage: fmt.Sprintf("[CALCULATION] Step %d, node name %s...", step, currentNodeName),
				})
				currentParams, err = em.triggerServerlessFunction(currentNode.FunctionName, currentParams, mqName)
				if err != nil {
					em.publishMessage(mqName, mqObject.MQmessage_Workflow{
						Isdone:        true,
						DataOrMessage: fmt.Sprintf("[FAIL] Step %d, node name %s...", step, currentNodeName),
					})
					return
				}
				currentNodeName = currentNode.Branch[0]
				currentNode = workflowObject.Spec.WorkflowNodes[currentNodeName]
				break
			}
		case workflow.SERVERLESS_NODETYPE_END:
			{
				step += 1
				em.publishMessage(mqName, mqObject.MQmessage_Workflow{
					Isdone:        false,
					DataOrMessage: fmt.Sprintf("[END] Step %d, node name %s...", step, currentNodeName),
				})
				currentParams, err = em.triggerServerlessFunction(currentNode.FunctionName, currentParams, mqName)
				if err != nil {
					em.publishMessage(mqName, mqObject.MQmessage_Workflow{
						Isdone:        true,
						DataOrMessage: fmt.Sprintf("[FAIL] Step %d, node name %s...", step, currentNodeName),
					})
					return
				}
				is_end = true
				break
			}
		}

		if is_end {
			break
		}
	}

	em.publishMessage(mqName, mqObject.MQmessage_Workflow{
		Isdone:        true,
		DataOrMessage: fmt.Sprintf("[FINISH] The workflow %s has finished.", workflowObject.Metadata.Name),
	})
}

func (em *EventsManager) publishMessage(routingKey string, msg mqObject.MQmessage_Workflow) {
	msgBody, err := json.Marshal(msg)
	if err != nil {
		println("Error marshalling msgBody")
		return
	}
	em.mqConn.Publish(
		minik8s_message.DefaultExchangeName,
		routingKey,
		"application/json",
		msgBody,
	)
}

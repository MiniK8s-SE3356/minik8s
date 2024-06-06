package formatprint

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/hpa"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/replicaset"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	"github.com/MiniK8s-SE3356/minik8s/pkg/gpu/types"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/function"
)

func PrintNodes(str string) {
	var nodes []struct {
		Metadata struct {
			ID     string            `json:"id"`
			Name   string            `json:"name"`
			Labels map[string]string `json:"labels"`
		} `json:"Metadata"`
		Status struct {
			Hostname   string   `json:"hostname"`
			IP         string   `json:"ip"`
			Condition  []string `json:"condition"`
			CPUPercent float64  `json:"cpuPercent"`
			MemPercent float64  `json:"memPercent"`
			NumPods    int      `json:"numPods"`
			UpdateTime string   `json:"updateTime"`
		} `json:"Status"`
	}

	err := json.Unmarshal([]byte(str), &nodes)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(writer, "ID\tName\tHostname\tIP\tCondition\tCPU Usage (%)\tMemory Usage (%)\tPod Count\tLast Update")

	for _, node := range nodes {
		// Join the condition array into a single string
		conditionStr := strings.Join(node.Status.Condition, ", ")
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%.2f\t%.2f\t%d\t%s\n",
			node.Metadata.ID, node.Metadata.Name, node.Status.Hostname,
			node.Status.IP, conditionStr, node.Status.CPUPercent*100, node.Status.MemPercent*100,
			node.Status.NumPods, node.Status.UpdateTime)
	}

	writer.Flush()
}

func PrintPods(str string) {
	var pods map[string]pod.Pod

	err := json.Unmarshal([]byte(str), &pods)
	if err != nil {
		fmt.Println("failed to unmarshal pods")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(writer, "ID\tName\tPhase\tPodIP\tNodeName\tCPU Usage (%)\tMemory Usage (%)")

	for _, pod := range pods {
		// Join the condition array into a single string
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%.6f\t%.6f\n",
			pod.Metadata.UUID, pod.Metadata.Name, pod.Status.Phase,
			pod.Status.PodIP, pod.Spec.NodeName, pod.Status.CPUUsage*100, pod.Status.MemoryUsage*100)
	}

	writer.Flush()
}

func PrintGPUJobs(str string) {
	var jobs map[string]types.SlurmJob

	err := json.Unmarshal([]byte(str), &jobs)
	if err != nil {
		fmt.Println("failed to unmarshal jobs")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(writer, "JobID\tName\tPartition\tState\tResult")

	for _, job := range jobs {
		// Join the condition array into a single string
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
			job.JobID, job.Metadata.Name, job.Partition, job.State, job.Result)
	}

	writer.Flush()
}

func PrintService(str string) {
	var services struct {
		ClusterIPArray []service.ClusterIP `json:"clusterIP"`
		NodePortArray  []service.NodePort  `json:"NodePort"`
	}

	err := json.Unmarshal([]byte(str), &services)
	if err != nil {
		fmt.Println("failed to unmarshal pods")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(writer, "ID\tName\tPhase\tIP\tType")

	for _, service := range services.ClusterIPArray {
		// Join the condition array into a single string
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
			service.Metadata.Id, service.Metadata.Name, service.Status.Phase,
			service.Metadata.Ip, service.Spec.Type)
	}

	for _, service := range services.NodePortArray {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
			service.Metadata.Id, service.Metadata.Name, service.Status.Phase,
			"null", service.Spec.Type)
	}

	writer.Flush()
}

func PrintReplicaset(str string) {
	var replicasets map[string]replicaset.Replicaset

	err := json.Unmarshal([]byte(str), &replicasets)
	if err != nil {
		fmt.Println("failed to unmarshal pods")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(writer, "ID\tName\tReady\tExpect")

	for _, rs := range replicasets {
		// Join the condition array into a single string
		fmt.Fprintf(writer, "%s\t%s\t%d\t%d\n",
			rs.Metadata.UUID, rs.Metadata.Name, rs.Status.ReadyReplicas,
			rs.Spec.Replicas)
	}

	writer.Flush()
}

func PrintHPA(str string) {
	var hpas map[string]hpa.HPA

	err := json.Unmarshal([]byte(str), &hpas)
	if err != nil {
		fmt.Println("failed to unmarshal pods")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(writer, "ID\tName\tReady\tMin\tMax")

	for _, h := range hpas {
		// Join the condition array into a single string
		fmt.Fprintf(writer, "%s\t%s\t%d\t%d\t%d\n",
			h.Metadata.UUID, h.Metadata.Name, h.Status.ReadyReplicas,
			h.Spec.MinReplicas, h.Spec.MaxReplicas)
	}

	writer.Flush()
}

func PrintFunction(str string) {

	var funcs []function.Function

	err := json.Unmarshal([]byte(str), &funcs)
	if err != nil {
		fmt.Println("failed to unmarshal func")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(writer, "ID\tName\tImageName")

	for _, h := range funcs {
		// Join the condition array into a single string
		fmt.Fprintf(writer, "%s\t%s\t%s\n",
			h.Metadata.UUID, h.Metadata.Name, h.Spec.ImageName)
	}

	writer.Flush()
}

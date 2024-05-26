package workflow

const (
	// 起始或中间计算节点
	SERVERLESS_NODETYPE_CALCULATION = "calculation"
	// 分支节点
	SERVERLESS_NODETYPE_BRANCH = "branch"
	// 终止节点
	SERVERLESS_NODETYPE_END = "end"
)

type Workflow struct {
	Kind       string           `json:"kind" yaml:"kind"` /*只允许Workflow*/
	ApiVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Metadata   WorkflowMetadata `json:"metadata" yaml:"metadata"`
	Spec       WorkflowSpec     `json:"spec" yaml:"spec"`
}

type WorkflowMetadata struct {
	Name string `json:"name" yaml:"name"`
}

type WorkflowSpec struct {
	// TODO：参数到底啥格式？
	Params string `json:"params" yaml:"params"`
	/* 入口点的Node名字 */
	EntryNodeName string `json:"entryNodeName" yaml:"entryNodeName"`
	/* nodeName->WorkflowNode */
	WorkflowNodes map[string]WorkflowNode `json:"workflowNodes" yaml:"workflowNodes"`
}

type WorkflowNode struct {
	/* 只能为SERVERLESS_NODETYPE_CALCULATION/BRANCH/END*/
	NodeType     string `json:"nodeType" yaml:"nodeType"`
	FunctionName string `json:"functionName" yaml:"functionName"` /* 函数名 */
	/* 当节点类型为calculation/end时，为空数组(尽量不要null)
	当节点类型为branch时，至少指定两个NodeName，用于下一轮的分支  */
	Branch []string `json:"branch" yaml:"branch"`
}

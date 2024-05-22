package message

type MsgType string

const (
	CreatePod MsgType = "create_pod"
	RemovePod MsgType = "remove_pod"
)

type Content interface{}

type Message struct {
	Type    MsgType `json:"type"`
	Content Content `json:"content"`
}

const (
	DefaultExchangeName       = "minik8s"
	DefaultSchedulerQueueName = "scheduler"
)

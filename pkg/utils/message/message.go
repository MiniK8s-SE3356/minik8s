package message

type MsgType string

const (
	CreatePod MsgType = "CreatePod"
)

type Content interface{}

type Message struct {
	Type    MsgType `json:"type"`
	Content Content `json:"content"`
}

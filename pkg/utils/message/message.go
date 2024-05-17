package message

type MsgType string

const (
	PodAdd MsgType = "PodAdd"
)

type Message struct {
	Type MsgType `json:"type"`
	Body string  `json:"body"`
}

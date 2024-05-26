package mqObject

type MQmessage_Workflow struct {
	Isdone        bool   `json:"isdone" yaml:"isdone"`
	DataOrMessage string `json:"dataormessage" yaml:"dataormessage"`
}

package httpobject

type HTTPResponse_callfunc struct {
	Result string `json:"result" yaml:"result"`
}

type HTTPRequest_callfunc struct {
	Params string `json:"params" yaml:"params"`
}

type HTTPResponse_GetAllServerlessFunction []string

type HTTPRequest_AddServerlessFuncPod struct {
	FuncName string `json:"funcName" yaml:"funcName"`
}

package httpobject

type HTTPResponse_callfunc struct {
	Data string `json:"data" yaml:"data"`
}

type HTTPRequest_callfunc struct {
	Params string `json:"params" yaml:"params"`
}

type HTTPResponse_GetAllServerlessFunction []string

type HTTPRequest_AddServerlessFuncPod struct {
	FuncName string `json:"funcName" yaml:"funcName"`
}

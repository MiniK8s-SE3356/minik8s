package config

const (
	HTTPURL_root_Template                     = "http://%s:%s"
	HTTPURL_callfunc_Template                 = "http://%s:5000/api/v1/callfunc"
	HTTPURL_AddServerlessFuncPod_Template     = "http://%s:%s/api/v1/AddServerlessFuncPod"
	HTTPURL_GetAllPod_Template                = "http://%s:%s/api/v1/GetAllPod"
	HTTPURL_GetAllServerlessFunction_Template = "http://%s:%s/api/v1/GetAllServerlessFunction"
)

var HTTPURL_AddServerlessFuncPod string
var HTTPURL_GetAllPod string
var HTTPURL_GetAllServerlessFunction string
var HTTPURL_root string

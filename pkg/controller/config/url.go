package config

const (
	HTTPURL_Template                     = "http://%s:%s"
	HTTPURL_GetAllService_Template       = "http://%s:%s/api/v1/GetAllService"
	HTTPURL_GetAllDNS_Template           = "http://%s:%s/api/v1/GetAllDNS"
	HTTPURL_UpdateDNS_Template           = "http://%s:%s/api/v1/UpdateDNS"
	HTTPURL_GetAllPod_Template           = "http://%s:%s/api/v1/GetAllPod"
	HTTPURL_UpdateService_Template       = "http://%s:%s/api/v1/UpdateService"
	HTTPURL_AddorDeleteEndpoint_Template = "http://%s:%s/api/v1/AddorDeleteEndpoint"
	HTTPURL_GetAllPersistVolume_Template = "http://%s:%s/api/v1/GetAllPersistVolume"
	HTTPURL_UpdatePersistVolume_Template = "http://%s:%s/api/v1/UpdatePersistVolume"
)

var HTTPURL string
var HTTPURL_GetAllService string
var HTTPURL_GetAllDNS string
var HTTPURL_UpdateDNS string
var HTTPURL_GetAllPod string
var HTTPURL_UpdateService string
var HTTPURL_AddorDeleteEndpoint string
var HTTPURL_GetAllPersistVolume string
var HTTPURL_UpdatePersistVolume string

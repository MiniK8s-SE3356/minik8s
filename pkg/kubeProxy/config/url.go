package config

const(
	HTTPURL_GetAllDNS_Template="http://%s:%s/api/v1/GetAllEndpoint"
	HTTPURL_GetAllService_Template="http://%s:%s/api/v1/GetAllService"
	HTTPURL_GetAllEndpoint_Template="http://%s:%s/api/v1/GetAllEndpoint"
) 

var HTTPURL_GetAllDNS string
var HTTPURL_GetAllEndpoint string
var HTTPURL_GetAllService string
var NGINX_IP string
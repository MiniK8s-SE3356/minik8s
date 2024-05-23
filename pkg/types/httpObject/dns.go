package httpobject

import "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/dns"

type HTTPResponse_GetAllDns []dns.Dns

type HTTPRequest_UpdateDns map[string]dns.Dns

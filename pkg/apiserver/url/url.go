package url

const (
	version = "v1"
	// 这里暂时是localhost
	rootURL     = "http://localhost:9000"
	prefix      = rootURL + "/api/" + version
	AddPodURL   = prefix + "/addPod"
	GetPodsURL  = prefix + "/getPods"
	GetNodesURL = prefix + "/getNodes"
)

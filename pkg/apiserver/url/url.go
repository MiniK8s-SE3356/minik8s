package url

const (
	version = "v1"
	// 这里暂时是localhost
	rootURL = "http://localhost:8080"
	prefix  = "/api/" + version

	AddPod      = prefix + "/AddPod"
	RemovePod   = prefix + "/RemovePod"
	GetPod      = prefix + "/GetPod"
	DescribePod = prefix + "/DescribePod"

	AddService         = prefix + "/AddService"
	RemoveService      = prefix + "/RemoveService"
	GetService         = prefix + "/GetService"
	GetAllService      = prefix + "/GetAllService"
	GetFilteredService = prefix + "/GetFilteredService"
	UpdateService      = prefix + "/UpdateService"
	DescribeService    = prefix + "/DescribeService"

	AddNode    = prefix + "/AddNode"
	GetNode    = prefix + "/GetNode"
	RemoveNode = prefix + "/RemoveNode"

	AddNamespace      = prefix + "/AddNamespace"
	RemoveNamespace   = prefix + "/RemoveNamespace"
	GetNamespace      = prefix + "/GetNamespace"
	DescribeNamespace = prefix + "/DescribeNamespace"

	AddReplicaset      = prefix + "/AddReplicaset"
	RemoveReplicaset   = prefix + "/RemoveReplicaset"
	GetReplicaset      = prefix + "/GetReplicaset"
	DescribeReplicaset = prefix + "/DescribeReplicaset"

	AddPodURL      = rootURL + prefix + "/AddPod"
	RemovePodURL   = rootURL + prefix + "/RemovePod"
	GetPodURL      = rootURL + prefix + "/GetPod"
	DescribePodURL = rootURL + prefix + "/DescribePod"

	AddServiceURL      = rootURL + prefix + "/AddService"
	RemoveServiceURL   = rootURL + prefix + "/RemoveService"
	GetServiceURL      = rootURL + prefix + "/GetService"
	DescribeServiceURL = rootURL + prefix + "/DescribeService"

	GetNodesURL = rootURL + prefix + "/GetNodes"

	AddNamespaceURL      = rootURL + prefix + "/AddNamespace"
	RemoveNamespaceURL   = rootURL + prefix + "/RemoveNamespace"
	GetNamespaceURL      = rootURL + prefix + "/GetNamespace"
	DescribeNamespaceURL = rootURL + prefix + "/DescribeNamespace"

	AddReplicasetURL      = rootURL + prefix + "/AddReplicaset"
	RemoveReplicasetURL   = rootURL + prefix + "/RemoveReplicaset"
	GetReplicasetURL      = rootURL + prefix + "/GetReplicaset"
	DescribeReplicasetURL = rootURL + prefix + "/DescribeReplicaset"
)

func GetURL(path string) string {
	return rootURL + path
}

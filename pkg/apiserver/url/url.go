package url

const (
	version = "v1"
	// 这里暂时是localhost
	rootURL = "http://localhost:8080"
	prefix  = "/api/" + version

	AddPod      = prefix + "/addPod"
	RemovePod   = prefix + "/removePod"
	GetPod      = prefix + "/getPod"
	DescribePod = prefix + "/describePod"

	AddService      = prefix + "/addService"
	RemoveService   = prefix + "/removeService"
	GetService      = prefix + "/getService"
	DescribeService = prefix + "/describeService"

	AddNode    = prefix + "/addNode"
	GetNode    = prefix + "/getNode"
	RemoveNode = prefix + "/removeNode"

	AddNamespace      = prefix + "/addNamespace"
	RemoveNamespace   = prefix + "/removeNamespace"
	GetNamespace      = prefix + "/getNamespace"
	DescribeNamespace = prefix + "/describeNamespace"

	AddReplicaset      = prefix + "/addReplicaset"
	RemoveReplicaset   = prefix + "/removeReplicaset"
	GetReplicaset      = prefix + "/getReplicaset"
	DescribeReplicaset = prefix + "/describeReplicaset"

	AddPodURL      = rootURL + prefix + "/addPod"
	RemovePodURL   = rootURL + prefix + "/removePod"
	GetPodURL      = rootURL + prefix + "/getPod"
	DescribePodURL = rootURL + prefix + "/describePod"

	AddServiceURL      = rootURL + prefix + "/addService"
	RemoveServiceURL   = rootURL + prefix + "/removeService"
	GetServiceURL      = rootURL + prefix + "/getService"
	DescribeServiceURL = rootURL + prefix + "/describeService"

	GetNodesURL = rootURL + prefix + "/getNodes"

	AddNamespaceURL      = rootURL + prefix + "/addNamespace"
	RemoveNamespaceURL   = rootURL + prefix + "/removeNamespace"
	GetNamespaceURL      = rootURL + prefix + "/getNamespace"
	DescribeNamespaceURL = rootURL + prefix + "/describeNamespace"

	AddReplicasetURL      = rootURL + prefix + "/addReplicaset"
	RemoveReplicasetURL   = rootURL + prefix + "/removeReplicaset"
	GetReplicasetURL      = rootURL + prefix + "/getReplicaset"
	DescribeReplicasetURL = rootURL + prefix + "/describeReplicaset"
)

func getURL(path string) string {
	return rootURL + path
}

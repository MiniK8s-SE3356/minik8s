package url

const (
	version = "v1"
	// 这里暂时是localhost
	rootURL = "http://localhost:9000"
	prefix  = rootURL + "/api/" + version

	AddPodURL      = prefix + "/addPod"
	RemovePodURL   = prefix + "/removePod"
	GetPodURL      = prefix + "/getPod"
	DescribePodURL = prefix + "/describePod"

	AddServiceURL      = prefix + "/addService"
	RemoveServiceURL   = prefix + "/removeService"
	GetServiceURL      = prefix + "/getService"
	DescribeServiceURL = prefix + "/describeService"

	GetNodesURL = prefix + "/getNodes"

	AddNamespaceURL      = prefix + "/addNamespace"
	RemoveNamespaceURL   = prefix + "/removeNamespace"
	GetNamespaceURL      = prefix + "/getNamespace"
	DescribeNamespaceURL = prefix + "/describeNamespace"

	AddReplicasetURL      = prefix + "/addReplicaset"
	RemoveReplicasetURL   = prefix + "/removeReplicaset"
	GetReplicasetURL      = prefix + "/getReplicaset"
	DescribeReplicasetURL = prefix + "/describeReplicaset"
)

package url

const (
	version = "v1"
	// 这里暂时是localhost
	RootURL = "http://localhost:8080"
	prefix  = "/api/" + version

	AddPod      = prefix + "/AddPod"
	RemovePod   = prefix + "/RemovePod"
	GetPod      = prefix + "/GetPod"
	GetAllPod   = prefix + "/GetAllPod"
	DescribePod = prefix + "/DescribePod"
	UpdatePod   = prefix + "/UpdatePod"

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

	NodeHeartBeat = prefix + "/NodeHeartBeat"

	AddNamespace      = prefix + "/AddNamespace"
	RemoveNamespace   = prefix + "/RemoveNamespace"
	GetNamespace      = prefix + "/GetNamespace"
	DescribeNamespace = prefix + "/DescribeNamespace"

	AddReplicaset      = prefix + "/AddReplicaset"
	RemoveReplicaset   = prefix + "/RemoveReplicaset"
	GetReplicaset      = prefix + "/GetReplicaset"
	DescribeReplicaset = prefix + "/DescribeReplicaset"

	GetAllEndpoint      = prefix + "/AddAllEndpoint"
	UpdateEndpointBatch = prefix + "/UpdateEndpointBatch"
	AddorDeleteEndpoint = prefix + "/AddorDeleteEndpoint"
)

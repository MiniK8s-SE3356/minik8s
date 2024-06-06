package url

var RootURL string = "http://localhost:8080"

const (
	version = "v1"
	prefix  = "/api/" + version

	AddPod               = prefix + "/AddPod"
	RemovePod            = prefix + "/RemovePod"
	GetPod               = prefix + "/GetPod"
	GetAllPod            = prefix + "/GetAllPod"
	UpdatePod            = prefix + "/UpdatePod"
	AddServerlessFuncPod = prefix + "/AddServerlessFuncPod"
	GetServerlessFuncPod = prefix + "/GetServerlessFuncPod"

	AddService         = prefix + "/AddService"
	RemoveService      = prefix + "/RemoveService"
	GetService         = prefix + "/GetService"
	GetAllService      = prefix + "/GetAllService"
	GetFilteredService = prefix + "/GetFilteredService"
	UpdateService      = prefix + "/UpdateService"

	AddNode    = prefix + "/AddNode"
	GetNode    = prefix + "/GetNode"
	RemoveNode = prefix + "/RemoveNode"

	NodeHeartBeat = prefix + "/NodeHeartBeat"

	AddNamespace    = prefix + "/AddNamespace"
	RemoveNamespace = prefix + "/RemoveNamespace"
	GetNamespace    = prefix + "/GetNamespace"

	AddReplicaset    = prefix + "/AddReplicaset"
	RemoveReplicaset = prefix + "/RemoveReplicaset"
	GetReplicaset    = prefix + "/GetReplicaset"

	AddHPA    = prefix + "/AddHPA"
	RemoveHPA = prefix + "/RemoveHPA"
	GetHPA    = prefix + "/GetHPA"
	UpdateHPA = prefix + "/updateHPA"

	AddDNS    = prefix + "/AddDNS"
	RemoveDNS = prefix + "/RemoveDNS"
	GetDNS    = prefix + "/GetDNS"
	UpdateDNS = prefix + "/UpdateDNS"
	GetAllDNS = prefix + "/GetAllDNS"

	GetAllEndpoint      = prefix + "/GetAllEndpoint"
	UpdateEndpointBatch = prefix + "/UpdateEndpointBatch"
	AddorDeleteEndpoint = prefix + "/AddorDeleteEndpoint"

	GetAllServerlessFunction = prefix + "/GetAllServerlessFunction"

	GetAllPersistVolume = prefix + "/GetAllPersistVolume"
	UpdatePersistVolume = prefix + "/UpdatePersistVolume"
	AddPV               = prefix + "/AddPV"
	AddPVC              = prefix + "/AddPVC"
	DeletePV			= prefix+"/DeletePV"
	DeletePVC			= prefix+"/DeletePVC"
)

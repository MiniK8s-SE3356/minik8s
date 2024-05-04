package ty

type ReplicaSetDesc struct {
	Name           string
	ReplicaNum     int
	ContainersDesc []ContainerDesc
}

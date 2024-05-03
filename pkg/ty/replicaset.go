package ty

// 这两个是从yaml读取的
type ContainerDesc struct {
	Name      string
	ImageName string
	Command   string
	port      int
}

type ReplicaSetDesc struct {
	Name           string
	ReplicaNum     int
	ContainersDesc []ContainerDesc
}

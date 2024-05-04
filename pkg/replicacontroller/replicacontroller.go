package replicacontroller

func handleOneReplicaSet() {
	// 从etcd获取当前所有pod，找出从属于当前replicaset的pod，数量不对就增加或删除
	// 然后写回etcd
}

func replicaControl() {
	//获取所有replicaset，再逐个调用handleOneReplicaSet处理
}

func Start() {
	go replicaControl()
}

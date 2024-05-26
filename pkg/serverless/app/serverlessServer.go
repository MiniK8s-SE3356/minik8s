package app

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/events"
)

type ServerlessServer struct {
	// 这里需要什么数据结构吗？
	events_manager (*events.EventsManager)
}

func NewServerlessServer() *ServerlessServer {
	fmt.Printf("New Serverless Server...\n")
	return &ServerlessServer{
		events_manager: events.NewEventsManager(),
	}
}

func (ss *ServerlessServer) Init() {
	fmt.Printf("Init Serverless Server...\n")
	ss.events_manager.Init()
	// TODO： 为Serving模块做初始化
	// TODO： 为Build模块做初始化
	// TODO： 为Event模块做初始化
}

func (ss *ServerlessServer) Run() {
	fmt.Printf("Run Serverless Server...\n")

	// TODO： 启动Serving，轮询执行routine,根据不同function的近期请求密度进行创建和回收

	// 启动Event.SyncRouteTableRoutine，定期更新内容

	// NOTE： Build和Event只有在调用时才会触发，不长期运行

	// TODO： 开启一些服务端口，用于接受来自kubectl的用户请求，并与api-server等组件进行交互
}

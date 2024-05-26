package app

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/events"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

type ServerlessServer struct {
	events_manager (*events.EventsManager)
}

func NewServerlessServer() *ServerlessServer {
	fmt.Printf("New Serverless Server...\n")
	return &ServerlessServer{
		events_manager: events.NewEventsManager(minik8s_message.DefaultMQConfig),
	}
}

func (ss *ServerlessServer) Init() {
	fmt.Printf("Init Serverless Server...\n")
	ss.events_manager.Init()
	// TODO： 为Serving模块做初始化
	// TODO： 为Build模块做初始化
}

func (ss *ServerlessServer) Run() {
	fmt.Printf("Run Serverless Server...\n")

	// 启动路由表的轮询更新机制
	go ss.events_manager.SyncRouteTableRoutine()

	// TODO： 启动Serving，轮询执行routine,根据不同function的近期请求密度进行创建和回收

	// NOTE： Build和Event只有在调用时才会触发，不长期运行

	// TODO： 开启一些服务端口，用于接受来自kubectl的用户请求，并与api-server等组件进行交互
}

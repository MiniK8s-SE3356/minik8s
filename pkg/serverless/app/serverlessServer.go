package app

import (
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/events"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/serving"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/server"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type ServerlessServer struct {
	events_manager (*events.EventsManager)
	s              (*server.Server)
}

func NewServerlessServer() *ServerlessServer {
	fmt.Printf("New Serverless Server...\n")
	return &ServerlessServer{
		events_manager: events.NewEventsManager(minik8s_message.DefaultMQConfig),
		s:              server.NewServer(),
	}
}

func (ss *ServerlessServer) Init() {
	fmt.Printf("Init Serverless Server...\n")
	ss.events_manager.Init()
	ss.s.Init()
	// TODO： 为Serving模块做初始化
	// TODO： 为Build模块做初始化
}

func (ss *ServerlessServer) Run() {
	fmt.Printf("Run Serverless Server...\n")

	// 启动路由表的轮询更新机制
	go ss.events_manager.SyncRouteTableRoutine()
	go poller.PollerStaticPeriod(10*time.Second, serving.ScaleFunctionPod, true)
	// TODO： 开启一些服务端口，用于接受来自kubectl的用户请求，并与api-server等组件进行交互
	ss.s.Run()
}

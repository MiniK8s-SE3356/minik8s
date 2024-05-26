package events

import (
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type EventsManager struct {
	route_table_manager *RouteTableManager
}

func NewEventsManager() *EventsManager {
	fmt.Printf("New EventsManager\n")
	return &EventsManager{
		route_table_manager: NewRouteTableManager(),
	}
}

func (em *EventsManager) Init() {
	fmt.Printf("Init EventsManager\n")
	em.route_table_manager.Init()
}

// func (em *EventsManager) Run() {
// 	fmt.Printf("Run EventsManager\n")
// }

func (em *EventsManager) SyncRouteTableRoutine() {
	// Never Return
	poller.PollerStaticPeriod(10*time.Second, em.route_table_manager.SyncRoutine, true)
}

package events

import "fmt"

type EventsManager struct{

}

func NewEventsManager() (*EventsManager) {
	fmt.Printf("New EventsManager\n")
	return &EventsManager{}
}

func (em *EventsManager)Init(){
	fmt.Printf("Init EventsManager\n")
}

func (em *EventsManager)Run(){
	fmt.Printf("Run EventsManager\n")
}
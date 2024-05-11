package servicesController

import "fmt"

type ServicesController struct {
}

func NewServicesController() *(ServicesController) {
	fmt.Printf("New ServicesController...\n")
	return &ServicesController{}
}

func (sc *ServicesController) Init() {
	fmt.Printf("Init ServicesController ...\n")

}

func (sc *ServicesController) Run() {
	fmt.Printf("Run ServicesController ...\n")

}

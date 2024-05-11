package endpointsController

import "fmt"

type EndpointsController struct {
	// NodePort

	// ClusterIP

	// Endpoints
}

func NewEndpointsController() *(EndpointsController) {
	fmt.Printf("New EndpointsController...\n")
	return &EndpointsController{}
}

func (ec *EndpointsController) Init() {
	fmt.Printf("Init EndpointsController ...\n")

}

func (ec *EndpointsController) Run() {
	fmt.Printf("Run EndpointsController ...\n")

}

func (ec *EndpointsController) SyncLocalData() {
	fmt.Printf("Sync LocalData ...\n")

}

func (ec *EndpointsController) RenewServiceStatus() {
	fmt.Printf("Renew Service Status ...\n")

}

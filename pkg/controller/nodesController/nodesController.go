package nodesController

import "fmt"

type NodesController struct {
}

func NewNodesController() *(NodesController) {
	fmt.Printf("New NodesController...\n")
	return &NodesController{}
}

func (ec *NodesController) Init() {
	fmt.Printf("Init NodesController ...\n")

}

func (ec *NodesController) Run() {
	fmt.Printf("Run NodesController ...\n")

}

package app

import (
	"fmt"
	"net/http"

	persistVolumeController "github.com/MiniK8s-SE3356/minik8s/pkg/controller/PersistVolumeController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/dnsController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/endpointsController"
	hpacontroller "github.com/MiniK8s-SE3356/minik8s/pkg/controller/hpaController"
	replicasetcontroller "github.com/MiniK8s-SE3356/minik8s/pkg/controller/replicasetController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/servicesController"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	endpointsController     *(endpointsController.EndpointsController)
	servicesController      *(servicesController.ServicesController)
	replicasetController    *(replicasetcontroller.ReplicasetController)
	hpaController           *(hpacontroller.HPAController)
	dnsController           *(dnsController.DnsController)
	persistvolumeController *(persistVolumeController.PersistVolumeController)
}

func NewController() *(Controller) {
	fmt.Printf("New Controller...\n")
	endpoints_controller := endpointsController.NewEndpointsController()
	services_controller := servicesController.NewServicesController()
	replicaset_controller := replicasetcontroller.NewReplicasetController()
	hpa_controller := hpacontroller.NewHPAController()
	dns_controller := dnsController.NewDnsController()
	persist_volume_controller := persistVolumeController.NewersistVolumeController()
	return &Controller{
		endpointsController:     endpoints_controller,
		servicesController:      services_controller,
		replicasetController:    replicaset_controller,
		hpaController:           hpa_controller,
		dnsController:           dns_controller,
		persistvolumeController: persist_volume_controller,
	}
}

func (co *Controller) Init() {
	fmt.Printf("Init Controller ...\n")
	co.endpointsController.Init()
	co.servicesController.Init()
	co.dnsController.Init()
	co.persistvolumeController.Init()
}

func (co *Controller) Run() {
	fmt.Printf("Run Controller ...\n")
	go co.endpointsController.Run()
	go co.servicesController.Run()
	go co.replicasetController.Run()
	go co.hpaController.Run()
	go co.dnsController.Run()
	go co.persistvolumeController.Run()
	// 打开一些监听端口
	r := gin.Default()
	co.bind(r)
	r.Run(":8082")
}

var url_AddPVImmediately = "/api/v1/AddPVImmediately"

func (co *Controller) bind(r *gin.Engine) {
	r.POST(url_AddPVImmediately, co.addPVImmediately)
}

func (co *Controller) addPVImmediately(c *gin.Context) {
	var requestMsg httpobject.HTTPReuqest_AddPVImmediately
	if err := c.ShouldBindJSON(&requestMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := co.persistvolumeController.CreatePVImmediately(requestMsg.PvName, requestMsg.PvType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "")
}

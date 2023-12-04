package controllers

import (
	"github.com/arslab/lwnsimulator/models"
	repo "github.com/arslab/lwnsimulator/repositories"

	gw "github.com/arslab/lwnsimulator/simulator/components/gateway"
	socketio "github.com/googollee/go-socket.io"
)

// SimulatorController interfaccia controller
type SimulatorController interface {
	Run() bool
	Stop() bool
	Status() bool
	GetIstance()
	AddWebSocket(*socketio.Conn)
	SaveBridgeAddress(models.AddressIP) error
	GetBridgeAddress() models.AddressIP
	GetGateways() []gw.Gateway
	AddGateway(*gw.Gateway) (int, int, error)
	UpdateGateway(*gw.Gateway) (int, error)
	DeleteGateway(int) bool
	ToggleStateGateway(int)
}

type simulatorController struct {
	repo      repo.SimulatorRepository
	onConnect func()
}

// NewSimulatorController return il controller
func NewSimulatorController(repo repo.SimulatorRepository) SimulatorController {
	return &simulatorController{
		repo: repo,
	}
}

func (c *simulatorController) GetIstance() {
	c.repo.GetIstance()
}

func (c *simulatorController) AddWebSocket(socket *socketio.Conn) {
	c.repo.AddWebSocket(socket)
}

func (c *simulatorController) Run() bool {
	return c.repo.Run()
}

func (c *simulatorController) Stop() bool {
	return c.repo.Stop()
}

func (c *simulatorController) Status() bool {
	return c.repo.Status()
}

func (c *simulatorController) SaveBridgeAddress(addr models.AddressIP) error {
	return c.repo.SaveBridgeAddress(addr)
}

func (c *simulatorController) GetBridgeAddress() models.AddressIP {
	return c.repo.GetBridgeAddress()
}

func (c *simulatorController) GetGateways() []gw.Gateway {
	return c.repo.GetGateways()
}

func (c *simulatorController) AddGateway(gateway *gw.Gateway) (int, int, error) {
	return c.repo.AddGateway(gateway)
}

func (c *simulatorController) UpdateGateway(gateway *gw.Gateway) (int, error) {
	return c.repo.UpdateGateway(gateway)
}

func (c *simulatorController) DeleteGateway(Id int) bool {
	return c.repo.DeleteGateway(Id)
}

func (c *simulatorController) ToggleStateGateway(Id int) {
	c.repo.ToggleStateGateway(Id)
}

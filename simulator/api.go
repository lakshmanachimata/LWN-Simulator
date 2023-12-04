package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/brocaar/lorawan"

	"github.com/arslab/lwnsimulator/codes"
	"github.com/arslab/lwnsimulator/models"

	f "github.com/arslab/lwnsimulator/simulator/components/forwarder"
	gw "github.com/arslab/lwnsimulator/simulator/components/gateway"
	c "github.com/arslab/lwnsimulator/simulator/console"
	"github.com/arslab/lwnsimulator/simulator/util"
	socketio "github.com/googollee/go-socket.io"
)

func GetIstance() *Simulator {

	var s Simulator

	s.State = util.Stopped

	s.loadData()

	s.ActiveGateways = make(map[int]int)

	s.Forwarder = *f.Setup()

	s.Console = c.Console{}

	return &s
}

func (s *Simulator) AddWebSocket(WebSocket *socketio.Conn) {
	s.Console.SetupWebSocket(WebSocket)
	s.Resources.AddWebSocket(WebSocket)
	s.SetupConsole()
}

func (s *Simulator) Run() {

	s.State = util.Running
	s.setup()

	s.Print("START", nil, util.PrintBoth)

	for _, id := range s.ActiveGateways {
		s.turnONGateway(id)
	}
}

func (s *Simulator) Stop() {

	s.State = util.Stopped
	s.Resources.ExitGroup.Add(len(s.ActiveGateways) - s.ComponentsInactiveTmp)

	for _, id := range s.ActiveGateways {
		s.Gateways[id].TurnOFF()
	}

	s.Resources.ExitGroup.Wait()

	s.saveStatus()

	s.Forwarder.Reset()

	s.Print("STOPPED", nil, util.PrintBoth)

	s.reset()

}

func (s *Simulator) SaveBridgeAddress(remoteAddr models.AddressIP) error {

	s.BridgeAddress = fmt.Sprintf("%v:%v", remoteAddr.Address, remoteAddr.Port)

	pathDir, err := util.GetPath()
	if err != nil {
		log.Fatal(err)
	}

	path := pathDir + "/simulator.json"

	bytes, err := json.MarshalIndent(&s, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = util.WriteConfigFile(path, bytes)
	if err != nil {
		log.Fatal(err)
	}

	s.Print("Gateway Bridge Address saved", nil, util.PrintOnlyConsole)

	return nil
}

func (s *Simulator) GetBridgeAddress() models.AddressIP {

	var rServer models.AddressIP
	if s.BridgeAddress == "" {
		return rServer
	}

	parts := strings.Split(s.BridgeAddress, ":")

	rServer.Address = parts[0]
	rServer.Port = parts[1]

	return rServer
}

func (s *Simulator) GetGateways() []gw.Gateway {

	var gateways []gw.Gateway

	for _, g := range s.Gateways {
		gateways = append(gateways, *g)
	}

	return gateways

}

func (s *Simulator) SetGateway(gateway *gw.Gateway, update bool) (int, int, error) {

	emptyAddr := lorawan.EUI64{0, 0, 0, 0, 0, 0, 0, 0}

	if gateway.Info.MACAddress == emptyAddr {

		s.Print("Error: MAC Address invalid", nil, util.PrintOnlyConsole)
		return codes.CodeErrorAddress, -1, errors.New("Error: MAC Address invalid")

	}

	if !update { //new

		gateway.Id = s.NextIDGw
		s.NextIDGw++

	} else {

		if s.Gateways[gateway.Id].IsOn() {
			return codes.CodeErrorDeviceActive, -1, errors.New("Gateway is running, unable update")
		}

	}

	code, err := s.searchName(gateway.Info.Name, gateway.Id, true)
	if err != nil {

		s.Print("Name already used", nil, util.PrintOnlyConsole)
		return code, -1, err

	}

	code, err = s.searchAddress(gateway.Info.MACAddress, gateway.Id, true)
	if err != nil {

		s.Print("DevEUI already used", nil, util.PrintOnlyConsole)
		return code, -1, err

	}

	if !gateway.Info.TypeGateway {

		if s.BridgeAddress == "" {
			return codes.CodeNoBridge, -1, errors.New("No gateway bridge configured")
		}

	}

	s.Gateways[gateway.Id] = gateway

	pathDir, err := util.GetPath()
	if err != nil {
		log.Fatal(err)
	}

	path := pathDir + "/gateways.json"
	s.saveComponent(path, &s.Gateways)
	path = pathDir + "/simulator.json"
	s.saveComponent(path, &s)

	s.Print("Gateway Saved", nil, util.PrintOnlyConsole)

	if gateway.Info.Active {

		s.ActiveGateways[gateway.Id] = gateway.Id

		if s.State == util.Running {
			s.Gateways[gateway.Id].Setup(&s.BridgeAddress, &s.Resources, &s.Forwarder)
			s.turnONGateway(gateway.Id)
		}

	} else {
		_, ok := s.ActiveGateways[gateway.Id]
		if ok {
			delete(s.ActiveGateways, gateway.Id)
		}
	}

	return codes.CodeOK, gateway.Id, nil
}

func (s *Simulator) DeleteGateway(Id int) bool {

	if s.Gateways[Id].IsOn() {
		return false
	}

	delete(s.Gateways, Id)
	delete(s.ActiveGateways, Id)

	pathDir, err := util.GetPath()
	if err != nil {
		log.Fatal(err)
	}

	path := pathDir + "/gateways.json"
	s.saveComponent(path, &s.Gateways)

	s.Print("Gateway Deleted", nil, util.PrintOnlyConsole)

	return true
}

func (s *Simulator) ToggleStateGateway(Id int) {

	if s.Gateways[Id].State == util.Stopped {
		s.turnONGateway(Id)
	} else {
		s.turnOFFGateway(Id)
	}

}

package forwarder

import (
	m "github.com/arslab/lwnsimulator/simulator/components/forwarder/models"
	"github.com/brocaar/lorawan"
)

func Setup() *Forwarder {

	f := Forwarder{
		Gateways: make(map[lorawan.EUI64]m.InfoGateway),
	}

	return &f

}

func (f *Forwarder) AddGateway(g m.InfoGateway) {

	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	f.Gateways[g.MACAddress] = g
}

func (f *Forwarder) DeleteGateway(g m.InfoGateway) {

	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	delete(f.Gateways, g.MACAddress)

}

func (f *Forwarder) Reset() {
	f = Setup()
}

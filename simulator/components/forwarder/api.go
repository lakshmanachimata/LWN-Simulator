package forwarder

import (
	dl "github.com/arslab/lwnsimulator/simulator/components/device/frames/downlink"
	m "github.com/arslab/lwnsimulator/simulator/components/forwarder/models"
	"github.com/arslab/lwnsimulator/simulator/resources/communication/buffer"
	pkt "github.com/arslab/lwnsimulator/simulator/resources/communication/packets"
	"github.com/brocaar/lorawan"
)

func Setup() *Forwarder {

	f := Forwarder{
		DevToGw:  make(map[lorawan.EUI64]map[lorawan.EUI64]*buffer.BufferUplink),            //1[devEUI] 2 [macAddress]
		GwtoDev:  make(map[uint32]map[lorawan.EUI64]map[lorawan.EUI64]*dl.ReceivedDownlink), //1[fre1] 2 [macAddress] 3[devEUI]
		Gateways: make(map[lorawan.EUI64]m.InfoGateway),
	}

	return &f

}

func (f *Forwarder) AddGateway(g m.InfoGateway) {

	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	f.Gateways[g.MACAddress] = g
}

func (f *Forwarder) DeleteDevice(DevEUI lorawan.EUI64) {

	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	for key := range f.DevToGw[DevEUI] {
		delete(f.DevToGw[DevEUI], key)
	}

	delete(f.DevToGw, DevEUI)

}

func (f *Forwarder) DeleteGateway(g m.InfoGateway) {

	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	delete(f.Gateways, g.MACAddress)

}

func (f *Forwarder) Uplink(data pkt.RXPK, DevEUI lorawan.EUI64) {

	f.Mutex.Lock()

	rxpk := createPacket(data)

	for _, up := range f.DevToGw[DevEUI] {
		up.Push(rxpk)
	}

	f.Mutex.Unlock()

}

func (f *Forwarder) Downlink(data *lorawan.PHYPayload, freq uint32, macAddress lorawan.EUI64) {

	f.Mutex.Lock()

	for _, dl := range f.GwtoDev[freq][macAddress] {
		dl.Push(data)
	}

	f.Mutex.Unlock()

}

func (f *Forwarder) Reset() {
	f = Setup()
}

package regional_parameters

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	c "github.com/arslab/lwnsimulator/simulator/components/device/features/channels"
	models "github.com/arslab/lwnsimulator/simulator/components/device/regional_parameters/models_rp"
	"github.com/brocaar/lorawan"
)

type Ql256 struct {
	Info models.Parameters
}

// manca un setup
func (eu *Ql256) Setup() {
	eu.Info.Code = Code_Ql256
	eu.Info.MinFrequency = 256100000
	eu.Info.MaxFrequency = 257500000
	eu.Info.FrequencyRX2 = 256700000
	eu.Info.DataRateRX2 = 7
	eu.Info.MinDataRate = 7
	eu.Info.MaxDataRate = 7
	eu.Info.MinRX1DROffset = 0
	eu.Info.MaxRX1DROffset = 0
	eu.Info.InfoGroupChannels = []models.InfoGroupChannels{
		{
			EnableUplink:       true,
			InitialFrequency:   256100000,
			OffsetFrequency:    200000,
			MinDataRate:        7,
			MaxDataRate:        7,
			NbReservedChannels: 3,
		},
	}
	eu.Info.InfoClassB.Setup(256700000, 256700000, 3, eu.Info.MinDataRate, eu.Info.MaxDataRate)

}

func (eu *Ql256) GetDataRate(datarate uint8) (string, string) {

	switch datarate {
	case 7:
		r := fmt.Sprintf("SF%vBW125", 12-datarate)
		return "LORA", r

	}
	return "", ""
}

func (eu *Ql256) FrequencySupported(frequency uint32) error {

	if frequency < eu.Info.MinFrequency || frequency > eu.Info.MaxFrequency {
		return errors.New("Frequency not supported")
	}

	return nil
}

func (eu *Ql256) DataRateSupported(datarate uint8) error {

	if datarate < eu.Info.MinDataRate || datarate > eu.Info.MaxDataRate {
		return errors.New("Invalid Data Rate")
	}

	return nil
}

func (eu *Ql256) GetCode() int {
	return Code_Ql256
}

func (eu *Ql256) GetChannels() []c.Channel {
	var channels []c.Channel

	for i := 0; i < eu.Info.InfoGroupChannels[0].NbReservedChannels; i++ {
		frequency := eu.Info.InfoGroupChannels[0].InitialFrequency + eu.Info.InfoGroupChannels[0].OffsetFrequency*uint32(i)
		ch := c.Channel{
			Active:            true,
			EnableUplink:      eu.Info.InfoGroupChannels[0].EnableUplink,
			FrequencyUplink:   frequency,
			FrequencyDownlink: frequency,
			MinDR:             0,
			MaxDR:             5,
		}
		channels = append(channels, ch)
	}

	return channels
}

func (eu *Ql256) GetMinDataRate() uint8 {
	return eu.Info.MinDataRate
}

func (eu *Ql256) GetMaxDataRate() uint8 {
	return eu.Info.MaxDataRate
}

func (eu *Ql256) GetNbReservedChannels() int {
	return eu.Info.InfoGroupChannels[0].NbReservedChannels
}

func (eu *Ql256) GetCodR(datarate uint8) string {
	return "4/5"
}

func (eu *Ql256) RX1DROffsetSupported(offset uint8) error {
	if offset >= eu.Info.MinRX1DROffset && offset <= eu.Info.MaxRX1DROffset {
		return nil
	}

	return errors.New("Invalid RX1DROffset")
}

func (eu *Ql256) LinkAdrReq(ChMaskCntl uint8, ChMask lorawan.ChMask,
	newDataRate uint8, channels *[]c.Channel) ([]bool, []error) {

	return linkADRReqForChannels(eu, ChMaskCntl, ChMask, newDataRate, channels)
}

func (eu *Ql256) SetupRX1(datarate uint8, rx1offset uint8, indexChannel int, dtime lorawan.DwellTime) (uint8, int) {

	DataRateRx1 := uint8(0)
	if datarate > rx1offset { //set data rate RX1
		DataRateRx1 = datarate - rx1offset
	}

	return DataRateRx1, indexChannel
}

func (eu *Ql256) SetupInfoRequest(indexChannel int) (string, int) {

	rand.Seed(time.Now().UTC().UnixNano())

	if indexChannel >= eu.GetNbReservedChannels() {
		indexChannel = rand.Int() % eu.GetNbReservedChannels()
	}

	_, datarate := eu.GetDataRate(5)
	return datarate, indexChannel
}

func (eu *Ql256) GetFrequencyBeacon() uint32 {
	return eu.Info.InfoClassB.FrequencyBeacon
}

func (eu *Ql256) GetDataRateBeacon() uint8 {
	return eu.Info.InfoClassB.DataRate
}

func (eu *Ql256) GetPayloadSize(datarate uint8, dTime lorawan.DwellTime) (int, int) {

	switch datarate {
	case 7:
		return 230, 222
	}
	return 0, 0
}

func (eu *Ql256) GetParameters() models.Parameters {
	return eu.Info
}

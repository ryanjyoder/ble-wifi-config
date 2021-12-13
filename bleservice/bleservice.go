package bleservice

import (
	"github.com/ryanjyoder/ble-wifi-config/wificonfig"
	"tinygo.org/x/bluetooth"
)

var (
	WifiSetupServiceUUID, _        = bluetooth.ParseUUID("a8232fe0-5b0f-11ec-bf63-0242ac130002")
	ConnectedToInternetCharUUID, _ = bluetooth.ParseUUID("a8232fe1-5b0f-11ec-bf63-0242ac130002")
	ConnectedSsidCharUUID, _       = bluetooth.ParseUUID("a8232fe2-5b0f-11ec-bf63-0242ac130002")
	SetSsidCharUUID, _             = bluetooth.ParseUUID("a8232fe3-5b0f-11ec-bf63-0242ac130002")
	SetSecretCharUUID, _           = bluetooth.ParseUUID("a8232fe4-5b0f-11ec-bf63-0242ac130002")
	ConnectionErrorCharUUID, _     = bluetooth.ParseUUID("a8232fe5-5b0f-11ec-bf63-0242ac130002")
)

type BleService struct {
	connectedToInternetChar bluetooth.Characteristic
	connectedSsidChar       bluetooth.Characteristic
	setSSidChar             bluetooth.Characteristic
	setSecretChar           bluetooth.Characteristic
	connectionErrorChar     bluetooth.Characteristic
	wifiConfigService       wificonfig.WifiConfigService
}

func (s *BleService) getServiceConfig() *bluetooth.Service {
	return &bluetooth.Service{
		UUID: WifiSetupServiceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &s.connectedToInternetChar,
				UUID:   ConnectedSsidCharUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
				Value:  make([]byte, 1),
			}, {
				Handle: &s.connectedSsidChar,
				UUID:   ConnectedSsidCharUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
				Value:  make([]byte, 33),
			}, {
				Handle:     &s.setSSidChar,
				UUID:       SetSecretCharUUID,
				Flags:      bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: s.setSsid,
				Value:      make([]byte, 33),
			}, {
				Handle:     &s.setSecretChar,
				UUID:       SetSecretCharUUID,
				Flags:      bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: s.setSsid,
				Value:      make([]byte, 65),
			}, {
				Handle: &s.connectionErrorChar,
				UUID:   ConnectionErrorCharUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
				Value:  make([]byte, 129),
			},
		},
	}
}

func (s *BleService) setSsid(client bluetooth.Connection, offset int, value []byte) {

}
func (s *BleService) setSecret(client bluetooth.Connection, offset int, value []byte) {

}

package bleservice

import (
	"fmt"

	"github.com/ryanjyoder/ble-wifi-config/connectivity"
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
	connectivityService connectivity.ConnectivityService
	wifiConfigService   wificonfig.WifiConfigService
	adapter             bluetooth.Adapter

	connectedToInternetChar bluetooth.Characteristic
	connectedSsidChar       bluetooth.Characteristic
	setSSidChar             bluetooth.Characteristic
	setSecretChar           bluetooth.Characteristic
	connectionErrorChar     bluetooth.Characteristic
}

func NewBleService(adapter bluetooth.Adapter, wifiConfigService wificonfig.WifiConfigService, connectivityService connectivity.ConnectivityService) (*BleService, error) {
	service := &BleService{
		wifiConfigService:   wifiConfigService,
		connectivityService: connectivityService,
		adapter:             adapter,
	}

	return service, nil
}

func (s *BleService) Start() error {
	err := s.adapter.Enable()
	if err != nil {
		return err
	}

	err = s.wireBleAttributes()
	if err != nil {
		return fmt.Errorf("error wiring ble attributes: %v", err)
	}

	adv := s.adapter.DefaultAdvertisement()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "WifiSetup",
		ServiceUUIDs: []bluetooth.UUID{WifiSetupServiceUUID},
	})
	if err != nil {
		return err
	}

	err = adv.Start()
	if err != nil {
		return err
	}
	return nil
}

func (s *BleService) wireBleAttributes() error {
	err := s.adapter.AddService(s.getServiceConfig())
	if err != nil {
		return err
	}

	s.connectivityService.NotifyInternetChange(s.handleInternetStatusChange)
	s.wifiConfigService.NotifySSIDChange(s.handleSsidChange)

	return nil
}

func (s *BleService) getServiceConfig() *bluetooth.Service {
	return &bluetooth.Service{
		UUID: WifiSetupServiceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &s.connectedToInternetChar,
				UUID:   ConnectedToInternetCharUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
				Value:  make([]byte, 1),
			}, {
				Handle: &s.connectedSsidChar,
				UUID:   ConnectedSsidCharUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
				Value:  make([]byte, 33),
			}, {
				Handle:     &s.setSSidChar,
				UUID:       SetSsidCharUUID,
				Flags:      bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: s.setSsid,
				Value:      make([]byte, 33),
			}, {
				Handle:     &s.setSecretChar,
				UUID:       SetSecretCharUUID,
				Flags:      bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: s.setSecret,
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
	fmt.Println("got write: ", string(value))
	ssid, err := bytesToStr(value)
	if err != nil {
		s.handleConnectionError(err)
		return
	}
	err = s.wifiConfigService.SetSSID(ssid)
	s.handleConnectionError(err)
}

func (s *BleService) setSecret(client bluetooth.Connection, offset int, value []byte) {
	secret, err := bytesToStr(value)
	if err != nil {
		s.handleConnectionError(err)
		return
	}
	err = s.wifiConfigService.SetSecret(secret)
	s.handleConnectionError(err)
}

func (s *BleService) handleSsidChange(ssid *string) {
	if ssid == nil {
		temp := ""
		ssid = &temp
	}

	strBytes, err := strTo33Bytes(*ssid)
	if err != nil {
		fmt.Println("error encoding string:", err)
		return
	}

	_, err = s.connectedSsidChar.Write(strBytes[:])
	if err != nil {
		fmt.Println("error writing to ssid char:", err)
		return
	}
}

func (s *BleService) handleInternetStatusChange(status bool) {
	statusArray := []byte{0x0}
	if status {
		statusArray[0] = 0x1
	}
	s.connectedToInternetChar.Write(statusArray)
}

func (s *BleService) handleConnectionError(err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	if errMsg == "" {
		fmt.Println("clearing error")
	}
	if len(errMsg) > 128 {
		errMsg = errMsg[:128]
	}
	errBytes, err := strTo129Bytes(errMsg)
	if err != nil {
		fmt.Println("Error serializing error msg:", err)
		return
	}

	_, err = s.connectionErrorChar.Write(errBytes[:])
	if err != nil {
		fmt.Println("error writing connection error:", err)
		return
	}
}

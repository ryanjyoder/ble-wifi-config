package wificonfig

import (
	"fmt"

	"github.com/Wifx/gonetworkmanager"
)

const (
	wirelessConnection           = "802-11-wireless"
	wirelessSecurity             = "802-11-wireless-security"
	wpaPsk                       = "wpa-psk"
	connectionSection            = "connection"
	connectionSectionID          = "id"
	connectionSectionAutoconnect = "autoconnect"
	ip4Section                   = "ipv4"
	ip4SectionMethod             = "method"
	ip4SectionNeverDefault       = "never-default"
	ip6Section                   = "ipv6"
	ip6SectionMethod             = "method"
	ipMethodIgnore               = "ignore"
	ipMethodAuto                 = "auto"

	connectionID = "ble-wifi-configured"
)

type WifiConfigService interface {
	// GetConnectedSSID can return nil,nil if the device is not connected to any wifi ssid
	GetConnectedSSID() (*string, error)
	NotifySSIDChange(SsidChangeCallback)
	SetSSID(ssid string) error
	SetSecret(secret string) error
}

type WifiConfigServiceNM struct {
	nm                 gonetworkmanager.NetworkManager
	settings           gonetworkmanager.Settings
	ssidChangeCallback SsidChangeCallback
	ssid               string
	secret             string
	connectionID       string
}

var _ WifiConfigService = &WifiConfigServiceNM{}

type SsidChangeCallback func(*string)

type NetworkManager interface {
}

func NewWifiConfigServiceNM(nm gonetworkmanager.NetworkManager, settings gonetworkmanager.Settings) (*WifiConfigServiceNM, error) {
	return &WifiConfigServiceNM{
		nm:           nm,
		settings:     settings,
		connectionID: connectionID,
	}, nil
}

func (s *WifiConfigServiceNM) GetConnectedSSID() (*string, error) {

	devices, err := s.nm.GetPropertyAllDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		ssid, err := connectedSSID(device)
		if err != nil {
			return nil, err
		}
		if ssid != nil {
			s.doCallback(ssid)
			return ssid, nil
		}
	}

	s.doCallback(nil)
	return nil, nil
}

func (s WifiConfigServiceNM) doCallback(ssid *string) {
	if s.ssidChangeCallback == nil {
		return
	}

	// make a copy of string
	if ssid != nil {
		temp := *ssid + ""
		ssid = &temp
	}

	go s.ssidChangeCallback(ssid)
}

func connectedSSID(device gonetworkmanager.Device) (*string, error) {
	deviceType, err := device.GetPropertyDeviceType()
	if err != nil {
		return nil, err
	}
	if deviceType != gonetworkmanager.NmDeviceTypeWifi {
		// not a wifi device, therefore no SSID
		return nil, nil
	}

	active, err := device.GetPropertyActiveConnection()
	if err != nil {
		return nil, err
	}
	if active == nil {
		return nil, nil
	}
	propConn, err := active.GetPropertyConnection()
	if err != nil {
		return nil, err
	}
	settings, err := propConn.GetSettings()
	if err != nil {
		return nil, err
	}

	if _, ok := settings[wirelessConnection]; !ok {
		return nil, nil
	}
	ssidI, ok := settings[wirelessConnection]["ssid"]
	if !ok {
		return nil, nil
	}
	ssidBytes, _ := ssidI.([]byte)
	ssid := string(ssidBytes)
	return &ssid, nil

}

func (s *WifiConfigServiceNM) NotifySSIDChange(callback SsidChangeCallback) {
	s.ssidChangeCallback = callback
}

func (s *WifiConfigServiceNM) SetSSID(ssid string) error {
	s.ssid = ssid
	if s.configsSet() {
		return s.configureWifi()
	}
	return nil
}
func (s *WifiConfigServiceNM) SetSecret(secret string) error {
	s.secret = secret
	if s.configsSet() {
		return s.configureWifi()
	}
	return nil
}

func (s WifiConfigServiceNM) configsSet() bool {
	return s.ssid != "" && s.secret != ""
}

func (s *WifiConfigServiceNM) configureWifi() error {
	newConfigs := getConfig(s.ssid, s.secret)

	currentSettings, err := s.checkForExistingConnection()
	if err != nil {
		return err
	}

	if currentSettings == nil {
		fmt.Println("creating new connection")
		_, err := s.settings.AddConnection(newConfigs)
		return err
	} else {
		fmt.Println("updating existing connection")
		return currentSettings.Update(newConfigs)
	}

	return nil
}

func getConfig(ssid string, secret string) map[string]map[string]interface{} {

	connection := make(map[string]map[string]interface{})
	connection[connectionSection] = make(map[string]interface{})
	connection[connectionSection][connectionSectionID] = connectionID
	connection[wirelessConnection] = make(map[string]interface{})
	connection[wirelessConnection]["ssid"] = []byte(ssid)
	connection[wirelessConnection]["security"] = wirelessSecurity
	connection[wirelessSecurity] = make(map[string]interface{})
	connection[wirelessSecurity]["key-mgmt"] = wpaPsk
	connection[wirelessSecurity]["psk"] = secret
	connection[connectionSection][connectionSectionAutoconnect] = true
	connection[ip4Section] = make(map[string]interface{})
	connection[ip4Section][ip4SectionMethod] = ipMethodAuto
	connection[ip4Section][ip4SectionNeverDefault] = true
	connection[ip6Section] = make(map[string]interface{})
	connection[ip6Section][ip6SectionMethod] = ipMethodIgnore

	return connection

}

func (s *WifiConfigServiceNM) checkForExistingConnection() (gonetworkmanager.Connection, error) {

	currentConnections, err := s.settings.ListConnections()
	if err != nil {
		return nil, err
	}

	for _, v := range currentConnections {
		connectionSettings, err := v.GetSettings()
		if err != nil {
			continue
		}

		currentConnectionSection := connectionSettings[connectionSection]
		if currentConnectionSection[connectionSectionID] == s.connectionID {
			return v, nil

		}
	}
	return nil, nil
}

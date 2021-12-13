package main

import (
	"fmt"
	"os"

	"github.com/Wifx/gonetworkmanager"
	"github.com/ryanjyoder/ble-wifi-config/wificonfig"
)

const (
	wirelessConnection           = "802-11-wireless"
	wirelessSecurity             = "802-11-wireless-security"
	wpaPsk                       = "wpa-psk"
	connectionSection            = "connection"
	connectionSectionID          = "id"
	connectionSectionAutoconnect = "autoconnect"
	ip4Section                   = "ipv4"
	ip4SectionAddress            = "address"
	ip4SectionPrefix             = "prefix"
	ip4SectionMethod             = "method"
	ip4SectionNeverDefault       = "never-default"
	ip6Section                   = "ipv6"
	ip6SectionMethod             = "method"
	ipMethodIgnore               = "ignore"
	ipMethodAuto                 = "ignore"
	desiredIP4Method             = "auto"
	desiredIP6Method             = "ignore"

	connectionID = "ble-wifi-configured"
)

func printVersion() error {
	/* Create new instance of gonetworkmanager */
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		return err
	}

	// Don't really need the network manager object per se
	// however knowing the version isn't bad
	var nmVersion string
	nmVersion, err = nm.GetPropertyVersion()
	if err != nil {
		return err
	}

	fmt.Println("Network Manager Version: " + nmVersion)
	return nil
}

func checkForExistingConnection() (gonetworkmanager.Connection, error) {
	// See if our connection already exists
	settings, err := gonetworkmanager.NewSettings()
	if err != nil {
		return nil, err
	}

	currentConnections, err := settings.ListConnections()
	if err != nil {
		return nil, err
	}

	for _, v := range currentConnections {
		connectionSettings, settingsError := v.GetSettings()
		if settingsError != nil {
			fmt.Println("settings error, continuing")
			continue
		}

		currentConnectionSection := connectionSettings[connectionSection]
		fmt.Println("ssid:", currentConnectionSection[connectionSectionID], connectionID)
		if currentConnectionSection[connectionSectionID] == connectionID {
			fmt.Println("Found setting")
			return v, nil

		}
	}
	return nil, nil
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

func createNewConnection(connection map[string]map[string]interface{}) error {

	settings, err := gonetworkmanager.NewSettings()

	if err != nil {
		return err
	}

	_, err = settings.AddConnection(connection)

	if err != nil {
		return err
	}
	return nil
}

func main() {

	nm, err := gonetworkmanager.NewNetworkManager()
	checkErr("error getting nm", err)

	settings, err := gonetworkmanager.NewSettings()
	checkErr("error getting settings", err)

	service, err := wificonfig.NewWifiConfigServiceNM(nm, settings)
	checkErr("error getting service", err)

	ssid, err := service.GetConnectedSSID()
	if err != nil {
		fmt.Println("error getting connected ssid:", err)
	}
	if ssid == nil {
		fmt.Println("not connected")
	} else {
		fmt.Println("connected ssid:", *ssid)
	}

	fmt.Println("setting ssid:", service.SetSSID(""))
	fmt.Println("secret:", service.SetSecret(""))

}

func checkErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg+": ", err)
		os.Exit(1)
	}
}

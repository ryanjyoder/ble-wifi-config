package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Wifx/gonetworkmanager"
	"github.com/ryanjyoder/ble-wifi-config/bleservice"
	"github.com/ryanjyoder/ble-wifi-config/connectivity"
	"github.com/ryanjyoder/ble-wifi-config/wificonfig"
	"tinygo.org/x/bluetooth"
)

func main() {

	nm, err := gonetworkmanager.NewNetworkManager()
	checkErr("error getting nm", err)

	settings, err := gonetworkmanager.NewSettings()
	checkErr("error getting settings", err)

	wifiService, err := wificonfig.NewWifiConfigServiceNM(nm, settings)
	checkErr("error getting service", err)

	connectivityService, err := connectivity.NewHttpConnectivityService(http.DefaultClient, "https://google.com")
	checkErr("eror getting connectivity service", err)

	service, err := bleservice.NewBleService(*bluetooth.DefaultAdapter, wifiService, connectivityService)
	checkErr("error getting bleservice", err)

	err = service.Start()
	checkErr("error starting bleservice", err)

	time.Sleep(15 * time.Minute)

	fmt.Println("shutdown after 15 minutes.")
}

func checkErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg+": ", err)
		os.Exit(1)
	}
}

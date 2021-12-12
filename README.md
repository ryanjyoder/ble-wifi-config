# ble-wifi-config
BLE Wifi Config

GATT Service UUID
a8232fe8-5b0f-11ec-bf63-0242ac130002

| Field      | Attribute Type | Length | Permissions | Description                                 | UUID                                 |
|------------|----------------|--------|-------------|---------------------------------------------|--------------------------------------|
| Wifi Setup | Service        | -      | -           | Service for configuring Wifi over Bluetooth | a8232fe0-5b0f-11ec-bf63-0242ac130002 |
| Connected to Internet | Characteristic | 1 | read/notify | Device currently connected to internet| a8232fe1-5b0f-11ec-bf63-0242ac130002 |
| Connected SSID | Characteristic | 33 |   | read/notify | The SSID currently connected            | a8232fe2-5b0f-11ec-bf63-0242ac130002 |
| Requested SSID | Characteristic | 33     | write   | The requested SSID to connect to            | a8232fe3-5b0f-11ec-bf63-0242ac130002 |
| Secret     | Characteristic | 65     | write       | Wifi passkey                                | a8232fe4-5b0f-11ec-bf63-0242ac130002 |
| Connection Error | Characteristic |  29 | read/notify | Last error when connecting to SSID       | a8232fe5-5b0f-11ec-bf63-0242ac130002 |
 
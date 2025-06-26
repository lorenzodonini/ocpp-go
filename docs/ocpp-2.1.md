## OCPP 2.1 Usage

OCPP 2.1 is now supported (experimental) in this library.

Requests and responses in OCPP 2.1 are handled the same way they were in v1.6 and 2.0.1.
There are additional message types on the websocket layer, which are handled automatically by the library.

The notable change is that there are now significantly more supported messages and profiles (feature sets),
which also require their own handlers to be implemented.

Below are very minimal setup code snippets, to get you started.
CSMS is now the equivalent of the Central System,
while the Charging Station is the new equivalent of a Charge Point.

Refer to the [examples folder](../example/2.1) for a full working example.
More in-depth documentation for v2.1 will follow.

### CSMS

To start a CSMS instance, run the following:

```go
import "github.com/lorenzodonini/ocpp-go/ocpp2.1"

csms := ocpp2.NewCSMS(nil, nil)

// Set callback handlers for connect/disconnect
csms.SetNewChargingStationHandler(func (chargingStation ocpp2.ChargingStationConnection) {
log.Printf("new charging station %v connected", chargingStation.ID())
})
csms.SetChargingStationDisconnectedHandler(func (chargingStation ocpp2.ChargingStationConnection) {
log.Printf("charging station %v disconnected", chargingStation.ID())
})

// Set handler for profile callbacks
handler := &CSMSHandler{}
csms.SetAuthorizationHandler(handler)
csms.SetAvailabilityHandler(handler)
csms.SetDiagnosticsHandler(handler)
csms.SetFirmwareHandler(handler)
csms.SetLocalAuthListHandler(handler)
csms.SetMeterHandler(handler)
csms.SetProvisioningHandler(handler)
csms.SetRemoteControlHandler(handler)
csms.SetReservationHandler(handler)
csms.SetTariffCostHandler(handler)
csms.SetTransactionsHandler(handler)


// Start central system
listenPort := 8887
log.Printf("starting CSMS")
csms.Start(listenPort, "/{ws}") // This call starts server in daemon mode and is blocking
log.Println("stopped CSMS")
```

#### Sending requests

Similarly to v2.0.1, you may send requests using the simplified API, e.g.

```go
err := csms.GetLocalListVersion(chargingStationID, myCallback)
if err != nil {
log.Printf("error sending message: %v", err)
}
```

Or you may build requests manually and send them using the asynchronous API.

#### Docker image

There is a Dockerfile and a docker image available upstream.

### Charging Station

To start a charging station instance, simply run the following:

```go
chargingStationID := "cs0001"
csmsUrl = "ws://localhost:8887"
chargingStation := ocpp2.NewChargingStation(chargingStationID, nil, nil)

// Set a handler for all callback functions
handler := &ChargingStationHandler{}
chargingStation.SetAvailabilityHandler(handler)
chargingStation.SetAuthorizationHandler(handler)
chargingStation.SetDataHandler(handler)
chargingStation.SetDiagnosticsHandler(handler)
chargingStation.SetDisplayHandler(handler)
chargingStation.SetFirmwareHandler(handler)
chargingStation.SetISO15118Handler(handler)
chargingStation.SetLocalAuthListHandler(handler)
chargingStation.SetProvisioningHandler(handler)
chargingStation.SetRemoteControlHandler(handler)
chargingStation.SetReservationHandler(handler)
chargingStation.SetSmartChargingHandler(handler)
chargingStation.SetTariffCostHandler(handler)
chargingStation.SetTransactionsHandler(handler)

// Connects to CSMS
err := chargingStation.Start(csmsUrl)
if err != nil {
log.Println(err)
} else {
log.Printf("connected to CSMS at %v", csmsUrl)
mainRoutine() // ... your program logic goes here
}

// Disconnect
chargingStation.Stop()
log.Println("disconnected from CSMS")
```

#### Sending requests

Similarly to v2.0.1 you may send requests using the simplified API (recommended), e.g.

```go
bootResp, err := chargingStation.BootNotification(provisioning.BootReasonPowerUp, "model1", "vendor1")
if err != nil {
log.Printf("error sending message: %v", err)
} else {
log.Printf("status: %v, interval: %v, current time: %v", bootResp.Status, bootResp.Interval, bootResp.CurrentTime.String())
}
```

Or you may build requests manually and send them using either the synchronous or asynchronous API.

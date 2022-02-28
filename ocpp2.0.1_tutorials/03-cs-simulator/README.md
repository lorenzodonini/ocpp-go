## Charging Station Simulator
The [previous example](../02-websocket-basic-auth) demonstrated how to test a CSMS Basic Auth connection using the `wscat` utility. This example shows how to use OCPP-Go to simulate a Charging Station using this programmatic [API](https://github.com/lorenzodonini/ocpp-go/blob/master/ocpp2.0.1/charging_station.go).

### OCPP-Go implementation
This example is very simple: the Charging Station connects to the CSMS using Basic Auth and immediatelly disconnects, without sending any OCPP messages (BootNotification, HeartBeat, etc) that will be covered by the next examples.
```
var log *logrus.Logger
var client *ws.Client
var cs ocpp2.ChargingStation

// Start function
func main() {
	// Connects to CSMS
	err := cs.Start(csUrl)
	if err != nil {
		log.Errorln(err)
	} else {
		log.Infof("Connected to CSMS at %v", csUrl)
		// Disconnect
		cs.Stop()
		log.Infof("Disconnected from CSMS")
	}
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	// Set this to DebugLevel if you want to retrieve verbose logs from the ocppj and websocket layers
	log.SetLevel(logrus.InfoLevel)

	ws.SetLogger(log.WithField("logger", "websocket"))

	client = ws.NewClient()
	client.SetBasicAuth("cs001", "s3cr3t")
	cs = ocpp2.NewChargingStation("cs001", nil, client)
}
```

### Running the example
#### Start the previous example CSMS
```
% cd 02-websockets-basic-auth
% go run .
INFO[2022-02-28T01:26:28-03:00] CSMS started on port 7777
INFO[2022-02-28T01:26:28-03:00] listening on tcp network :7777                logger=websocket
```
#### Start the Charging Station
```
% cd 03-cs-simulator
% go run .
INFO[2022-02-28T01:26:41-03:00] connecting to server                          logger=websocket
INFO[2022-02-28T01:26:41-03:00] connected to server as cs001                  logger=websocket
INFO[2022-02-28T01:26:41-03:00] Connected to CSMS at ws://localhost:7777
INFO[2022-02-28T01:26:41-03:00] closing connection to server                  logger=websocket
INFO[2022-02-28T01:26:41-03:00] Disconnected from CSMS
```

### CSMS Authentication Logs
```
INFO[2022-02-28T01:26:41-03:00] Charging Station connected                    client=cs001
INFO[2022-02-28T01:26:41-03:00] closed connection to cs001                    logger=websocket
INFO[2022-02-28T01:26:41-03:00] Charging Station disconnected                 client=cs001
```
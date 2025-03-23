### WebSocket Communication
The Charging Station Management System (CSMS) acts as a WebSocket server and the Charging Station (CS) acts as a WebSocket client. OCPP imposes extra constraints on the URL and the WebSocket subprotocol. If the CSMS endpoint URL is `ws://localhost:7777`:

- Each CS will connect to the server by appending its own identifier to the URL. Ex.:
  - `ws://localhost:7777/cs123            # ID: cs123`
  - `ws://localhost:7777/space%20escaped  # ID: space escaped`
- The CS MUST send a `Sec-WebSocket-Protocol` header:
  - `Sec-WebSocket-Protocol: ocpp2.0.1, ocpp1.6`

The example below demonstrates the connection between the CSMS and a CS. It instantiates an OCPP 2.0.1 server that will print log messages when a Charging Station connects and disconnects.

```
var csms = ocpp2.NewCSMS(nil, nil)

// Configure the CSMS
func init() {
	csms.SetNewChargingStationHandler(func(cs ocpp2.ChargingStationConnection) {
		fmt.Printf("Charging Station connected: %v\n", cs.ID())
	})
	csms.SetChargingStationDisconnectedHandler(func(cs ocpp2.ChargingStationConnection) {
		fmt.Printf("Charging Station disconnected: %v\n", cs.ID())
	})
}

// Start the CSMS
func main() {
	fmt.Println("Starting CSMS on port 7777")
	csms.Start(7777, "/{charging_station_id}")
}
```

Let's run the example:
```
% cd ocpp2.0.1_tutorials/01-websocket-connection
% go run .
INFO[2022-02-27T20:20:06-03:00] CSMS started on port 7777
```

And use `wscat` to simulate a Charging Station client connection in another tab:
```
% npm install wscat

# If the protocol header is not sent,
# the CSMS server will automatically disconnect.
% wscat -c ws://localhost:7777/123
Connected (press CTRL+C to quit)
Disconnected (code: 1002, reason: "invalid or unsupported subprotocol")

# If the connection is established,
# it can be used to send messages to the CSMS.
% wscat -s ocpp2.0.1 -c ws://localhost:7777/123
Connected (press CTRL+C to quit)
> Type something and press enter

# The connection can be closed manually (Control+C)
# or automatically after a timeout.

# Server logs
INFO[2022-02-27T20:20:50-03:00] Charging Station connected      client=123
INFO[2022-02-27T20:20:59-03:00] Charging Station disconnected   client=123
```
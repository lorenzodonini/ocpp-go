## Handler Methods

The [example 02](../02-websocket-basic-auth) demonstrated how to setup some basic handlers to deal with a Charging Station connection/disconnection and authetication.

A way of organizing the code consists in writing them as handler methods. Indeed this is how the code is organized in the [OCPP-Go 1.6 example](https://github.com/lorenzodonini/ocpp-go/blob/1c65522f8ab806fc626f04adbccdebbaed186f2a/example/1.6/cs/handler.go#L64):

https://github.com/lorenzodonini/ocpp-go/blob/1c65522f8ab806fc626f04adbccdebbaed186f2a/example/1.6/cs/handler.go#L64-L73

Let's refactor the following code and get familiar with the OCPP-Go internals:

```
import (
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ws"
)

var server *ws.Server
var csms ocpp2.CSMS

...

server.SetBasicAuthHandler(func(user string, pass string) bool {
  ok := authenticate(user, pass) // ... check for user and pass correctness
  return ok
})

csms = ocpp2.NewCSMS(nil, server)
csms.SetNewChargingStationHandler(func(cs ocpp2.ChargingStationConnection) {
  log.WithField("client", cs.ID()).Info("Charging Station connected")
})
csms.SetChargingStationDisconnectedHandler(func(cs ocpp2.ChargingStationConnection) {
  log.WithField("client", cs.ID()).Info("Charging Station disconnected")
})
```

The basic authentication is part of the interface of the WesSocket server: https://github.com/lorenzodonini/ocpp-go/blob/1c65522f8ab806fc626f04adbccdebbaed186f2a/ws/websocket.go#L227-L230

It's a function that receives the arguments username and password, then returns true if the authentication succeeded.

```
type CsmsHandler struct {
}

func (csms *CsmsHandler) Authenticate(username, password string) bool {
	ok := username == "cs0001" && password == "s3cr3t"
	log.Info("Authenticated", username, "?", ok)
	return ok
}
```

Having the `CsmsHandler` type as the [receiver](https://go.dev/tour/methods/1) of the methos Authenticate, the setup becomes a bit simpler:

```
handler := &CsmsHandler{}
server.SetBasicAuthHandler(handler.Authenticate)
```

Now, take a look at the Charging Station connection handling interface:

https://github.com/lorenzodonini/ocpp-go/blob/1c65522f8ab806fc626f04adbccdebbaed186f2a/ocpp2.0.1/v2.go#L33-L39

https://github.com/lorenzodonini/ocpp-go/blob/1c65522f8ab806fc626f04adbccdebbaed186f2a/ocpp2.0.1/v2.go#L370-L373

The CSMS interface defines these 2 connection handler methods, which argument is a `ChargingStationConnectionHandler` - a function that receives as argument a reference to a `ChargingStationConnection`. Let's expand our `CsmsHandler`:

```
type CsmsHandler struct {
}

func (csms *CsmsHandler) Authenticate(username, password string) bool {
	return username == 'cs0001' && password == 's3cr3t'
}

func (csms *CsmsHandler) OnChargingStationConnection(cs ChargingStationConnection) {
  log.WithField("client", cs.ID()).Info("Charging Station connected")
}

func (csms *CsmsHandler) OnChargingStationDisconnection(cs ChargingStationConnection) {
  log.WithField("client", cs.ID()).Info("Charging Station disconnected")
}
```

The resulting setup becomes:

```
handler := &CsmsHandler{}
server.SetBasicAuthHandler(handler.Authenticate)
csms.SetNewChargingStationHandler(handler.OnChargingStationConnection)
csms.SetChargingStationDisconnectedHandler(handler.OnChargingStationDisconnection)
```

## Conclusion
Although this is a refactoring tutorial, it demonstrates how to navigate through the OCPP-Go repository and get familiar with the types, interfaces and good coding practices.
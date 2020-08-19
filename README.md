# ocpp-go

[![Build Status](https://travis-ci.org/lorenzodonini/ocpp-go.svg?branch=master)](https://travis-ci.org/lorenzodonini/ocpp-go)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://godoc.org/github.com/lorenzodonini/ocpp-go)
[![Coverage Status](https://coveralls.io/repos/github/lorenzodonini/ocpp-go/badge.svg?branch=master)](https://coveralls.io/github/lorenzodonini/ocpp-go?branch=master)
[![Go report](https://goreportcard.com/badge/github.com/lorenzodonini/ocpp-go)](https://goreportcard.com/report/github.com/lorenzodonini/ocpp-go)

Open Charge Point Protocol implementation in Go.

The library targets modern charge points and central systems, running OCPP version 1.6+.

Given that SOAP will no longer be supported in future versions of OCPP, only OCPP-J is supported in this library.
There are currently no plans of supporting OCPP-S.

## Status & Roadmap

**Note: Releases 0.10.0 introduced breaking changes in some API, due to refactoring. The functionality remains the same, but naming changed.**

Planned milestones and features:

- [x] OCPP 1.6
- [ ] OCPP 2.0 

## OCPP 1.6 Usage

Go version 1.11+ is required.

```sh
go get github.com/lorenzodonini/ocpp-go
```

You will also need to fetch some dependencies:
```sh
cd <path-to-ocpp-go>
export GO111MODULE=on
go mod download
```

Your application may either act as a [Central System](#central-system) (server) or as a [Charge Point](#charge-point) (client).

### Central System

If you want to integrate the library into your custom Central System, you must implement the callbacks defined in the profile interfaces, e.g.:
```go
import (
    "github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
    "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
    "time"
)

const defaultHeartbeatInterval = 600

type CentralSystemHandler struct {
	// ... your own state variables
}

func (handler *CentralSystemHandler) OnAuthorize(chargePointId string, request *core.AuthorizeRequest) (confirmation *core.AuthorizeConfirmation, err error) {
	// ... your own custom logic
	return core.NewAuthorizationConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted)), nil
}

func (handler *CentralSystemHandler) OnBootNotification(chargePointId string, request *core.BootNotificationRequest) (confirmation *core.BootNotificationConfirmation, err error) {
	// ... your own custom logic
	return core.NewBootNotificationConfirmation(types.NewDateTime(time.Now()), defaultHeartbeatInterval, core.RegistrationStatusAccepted), nil
}

// further callbacks... 
```

Every time a request from the charge point comes in, the respective callback function is called.
For every callback you must return either a confirmation or an error. The result will be sent back automatically to the charge point.
The callback is invoked inside a dedicated goroutine, so you don't have to worry about synchronization.

You need to implement at least all other callbacks defined in the `core.CentralSystemHandler` interface.

Depending on which OCPP profiles you want to support in your application, you will need to implement additional callbacks as well.

To start a central system instance, simply run the following:
```go
centralSystem := ocpp16.NewCentralSystem(nil, nil)

// Set callback handlers for connect/disconnect
centralSystem.SetNewChargePointHandler(func(chargePointId string) {
	log.Printf("new charge point %v connected", chargePointId)
})
centralSystem.SetChargePointDisconnectedHandler(func(chargePointId string) {
	log.Printf("charge point %v disconnected", chargePointId)
})

// Set handler for profile callbacks
handler := &CentralSystemHandler{}
centralSystem.SetCoreHandler(handler)

// Start central system
listenPort := 8887
log.Printf("starting central system")
centralSystem.Start(listenPort, "/{ws}") // This call starts server in daemon mode and is blocking
log.Println("stopped central system")
```

#### Sending requests

To send requests to the charge point, you may either use the simplified API:
```go
err := centralSystem.ChangeAvailability("1234", myCallback, 1, core.AvailabilityTypeInoperative)
if err != nil {
	log.Printf("error sending message: %v", err)
}
```

or create a message manually:
```go
request := core.NewChangeAvailabilityRequest(1, core.AvailabilityTypeInoperative)
err := centralSystem.SendRequestAsync("clientId", request, callbackFunction)
if err != nil {
	log.Printf("error sending message: %v", err)
}
```

In both cases, the request is sent asynchronously and the function returns right away. 
You need to write the callback function to check for errors and handle the confirmation on your own:
```go
myCallback := func(confirmation *core.ChangeAvailabilityConfirmation, e error) {
	if e != nil {
		log.Printf("operation failed: %v", e)
	} else {
		log.Printf("status: %v", confirmation.Status)
		// ... your own custom logic
	}
}
```

Since the initial `centralSystem.Start` call blocks forever, you may want to wrap it in a goroutine (that is, if you need to run other operations on the main thread).

#### Example

You can take a look at the [full example](./example/1.6/cs/central_system_sim.go).
To run it, simply execute:
```bash
go run ./example/1.6/cs/*.go
```

#### Docker

A containerized version of the central system example is available:
```bash
docker pull ldonini/ocpp1.6-central-system:latest
docker run -it -p 8887:8887 --rm --name central-system ldonini/ocpp1.6-central-system:latest
```

You can also run it directly using docker-compose:
```sh
docker-compose -f example/1.6/docker-compose.yml up central-system
```

#### TLS

If you wish to test the central system using TLS, make sure you put your self-signed certificates inside the `example/1.6/certs` folder.

Feel free to use the utility script `cd example/1.6 && ./create-test-certificates.sh` for generating test certificates. 

Then run the following:
```
docker-compose -f example/1.6/docker-compose.tls.yml up central-system
```

### Charge Point

If you want to integrate the library into your custom Charge Point, you must implement the callbacks defined in the profile interfaces, e.g.:
```go
import (
    "github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
    "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

type ChargePointHandler struct {
	// ... your own state variables
}

func (handler *ChargePointHandler) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	// ... your own custom logic
	return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusAccepted), nil
}

func (handler *ChargePointHandler) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	// ... your own custom logic
	return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusAccepted), nil
}

// further callbacks...
```

When a request from the central system comes in, the respective callback function gets invoked.
For every callback you must return either a confirmation or an error. The result will be sent back automatically to the central system.

You need to implement at least all other callbacks defined in the `core.ChargePointHandler` interface.

Depending on which OCPP profiles you want to support in your application, you will need to implement additional callbacks as well.

To start a charge point instance, simply run the following:
```go
chargePointId := "cp0001"
csUrl = "ws://localhost:8887"
chargePoint := ocpp16.NewChargePoint(chargePointId, nil, nil)

// Set a handler for all callback functions
handler := &ChargePointHandler{}
chargePoint.SetCoreHandler(handler)

// Connects to central system
err := chargePoint.Start(csUrl)
if err != nil {
	log.Println(err)
} else {
	log.Printf("connected to central system at %v", csUrl) 
	mainRoutine() // ... your program logic goes here
}

// Disconnect
chargePoint.Stop()
log.Printf("disconnected from central system")
```

#### Sending requests

To send requests to the central station, you have two options. You may either use the simplified synchronous blocking API (recommended):
```go
bootConf, err := chargePoint.BootNotification("model1", "vendor1")
if err != nil {
	log.Fatal(err)
}
log.Printf("status: %v, interval: %v, current time: %v", bootConf.Status, bootConf.Interval, bootConf.CurrentTime.String())
// ... do something with the confirmation
```

or create a message manually:
```go
request := core.NewBootNotificationRequest("model1", "vendor1")
```

You can then decide to send the message using a synchronous blocking call:
```go
// Synchronous call
confirmation, err := chargePoint.SendRequest(request)
if err != nil {
	log.Printf("error sending message: %v", err)
}
bootConf := confirmation.(*core.BootNotificationConfirmation)
// ... do something with the confirmation
```
or an asynchronous call:
```go
// Asynchronous call
err := chargePoint.SendRequestAsync(request, callbackFunction)
if err != nil {
	log.Printf("error sending message: %v", err)
}
```

In the latter case, you need to write the callback function and check for errors on your own:
```go
callback := func(confirmation ocpp.Response, e error) {
	bootConf := confirmation.(*core.BootNotificationConfirmation)
	if e != nil {
		log.Printf("operation failed: %v", e)
	} else {
		log.Printf("status: %v", bootConf.Status)
		// ... your own custom logic
	}
}
```

When creating a message manually, you always need to perform type assertion yourself, as the `SendRequest` and `SendRequestAsync` APIs use generic `Request` and `Confirmation` interfaces.

#### Example
You can take a look at the [full example](./example/1.6/cp/charge_point_sim.go).
To run it, simply execute:
```bash
CLIENT_ID=chargePointSim CENTRAL_SYSTEM_URL=ws://<host>:8887 go run example/1.6/cp/*.go
```

You need to specify the URL of a running central station server via environment variable, so the charge point can reach it.

#### Docker

A containerized version of the charge point example is available:
```bash
docker pull ldonini/ocpp1.6-charge-point:latest
docker run -e CLIENT_ID=chargePointSim -e CENTRAL_SYSTEM_URL=ws://<host>:8887 -it --rm --name charge-point ldonini/ocpp1.6-charge-point:latest
```

You need to specify the host, on which the central system is running, in order for the charge point to connect to it.

You can also run it directly using docker-compose:
```sh
docker-compose -f example/1.6/docker-compose.yml up charge-point
```

#### TLS

If you wish to test the charge point using TLS, make sure you put your self-signed certificates inside the `example/1.6/certs` folder.

Feel free to use the utility script `cd example/1.6 && ./create-test-certificates.sh` for generating test certificates. 

Then run the following:
```
docker-compose -f example/1.6/docker-compose.tls.yml up charge-point
```

## OCPP 2.0 Usage

Documentation will follow, once the protocol is fully implemented.

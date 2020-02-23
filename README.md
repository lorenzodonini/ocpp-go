# ocpp-go

[![Build Status](https://travis-ci.org/lorenzodonini/ocpp-go.svg?branch=master)](https://travis-ci.org/lorenzodonini/ocpp-go)

Open Charge Point Protocol implementation in Go.

The library targets modern charge points and central systems, running OCPP version 1.6+.

Given that SOAP will no longer be supported in future versions of OCPP, only OCPP-J is supported in this library.
There are currently no plans of supporting OCPP-S.

## Roadmap

Planned milestones and features:

- [x] OCPP 1.6
- [ ] OCPP 2.0

**Note: The library is still a WIP, therefore expect some APIs to change.** 

## Usage

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

Your application may either act as a Central System (server) or as a Charge Point (client).

### Central System

If you want to integrate the library into your custom Central System, you must implement the callbacks defined in the profile interfaces, e.g.:
```go
type CentralSystemHandler struct {
	// ... your own state variables
}

func (handler * CentralSystemHandler) OnAuthorize(chargePointId string, request *ocpp16.AuthorizeRequest) (confirmation *ocpp16.AuthorizeConfirmation, err error) {
	// ... your own custom logic
	return ocpp16.NewAuthorizationConfirmation(ocpp16.NewIdTagInfo(ocpp16.AuthorizationStatusAccepted)), nil
}

func (handler * CentralSystemHandler) OnBootNotification(chargePointId string, request *ocpp16.BootNotificationRequest) (confirmation *ocpp16.BootNotificationConfirmation, err error) {
	// ... your own custom logic
	return ocpp16.NewBootNotificationConfirmation(ocpp16.NewDateTime(time.Now()), defaultHeartbeatInterval, ocpp16.RegistrationStatusAccepted), nil
}

// further callbacks... 
```

Every time a request from the charge point comes in, the respective callback function is called.
For every callback you must return either a confirmation or an error. The result will be sent back automatically to the charge point.
The callback is invoked inside a dedicated goroutine, so you don't have to worry about synchronization.

You need to implement at least all other callbacks defined in the `ocpp16.CentralSystemCoreListener` interface.

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
centralSystem.SetCentralSystemCoreListener(handler)

// Start central system
listenPort := 8887
log.Printf("starting central system")
centralSystem.Start(listenPort, "/{ws}") // This call starts server in daemon mode and is blocking
log.Println("stopped central system")
```

#### Sending requests

To send requests to the charge point, you may either use the simplified API:
```go
err := centralSystem.ChangeAvailability("1234", myCallback, 1, ocpp16.AvailabilityTypeInoperative)
if err != nil {
	log.Printf("error sending message: %v", err)
}
```

or create a message manually:
```go
request := ocpp16.NewChangeAvailabilityRequest(1, ocpp16.AvailabilityTypeInoperative)
err := centralSystem.SendRequestAsync("clientId", request, callbackFunction)
if err != nil {
	log.Printf("error sending message: %v", err)
}
```

In both cases, the request is sent asynchronously and the function returns right away. 
You need to write the callback function to check for errors and handle the confirmation on your own:
```go
myCallback := func(confirmation *ocpp16.ChangeAvailabilityConfirmation, e error) {
	if e != nil {
		log.Printf("operation failed: %v", e)
	} else {
		log.Printf("status: %v", confirmation.Status)
		// ... your own custom logic
	}
}
```

Since the initial `centralSystem.Start` call blocks forever, you may want to wrap it in a goroutine (that is, if you need to send requests to charge points form the main thread).

#### Example

You can take a look at the full example inside `central_system_sim.go`.
To run it, simply execute:
```bash
go run ./example/cs/central_system_sim.go
```

#### Docker

A containerized version of the central system example is available:
```bash
docker pull ldonini/ocpp1.6-central-system:latest
docker run -it -p 8887:8887 --rm --name central-system ldonini/ocpp1.6-central-system:latest
```

You can also build the docker image from source, using:
```sh
docker-compose up central_system
```

### Charge Point

If you want to integrate the library into your custom Charge Point, you must implement the callbacks defined in the profile interfaces, e.g.:
```go
type ChargePointHandler struct {
	// ... your own state variables
}

func (handler * ChargePointHandler) OnChangeAvailability(request *ocpp16.ChangeAvailabilityRequest) (confirmation *ocpp16.ChangeAvailabilityConfirmation, err error) {
	// ... your own custom logic
	return ocpp16.NewChangeAvailabilityConfirmation(ocpp16.AvailabilityStatusAccepted), nil
}

func (handler * ChargePointHandler) OnChangeConfiguration(request *ocpp16.ChangeConfigurationRequest) (confirmation *ocpp16.ChangeConfigurationConfirmation, err error) {
	// ... your own custom logic
	return ocpp16.NewChangeConfigurationConfirmation(ocpp16.ConfigurationStatusAccepted), nil
}

// further callbacks...
```

When a request from the central system comes in, the respective callback function is called.
For every callback you must return either a confirmation or an error. The result will be sent back automatically to the central system.

You need to implement at least all other callbacks defined in the `ocpp16.ChargePointCoreListener` interface.

Depending on which OCPP profiles you want to support in your application, you will need to implement additional callbacks as well.

To start a charge point instance, simply run the following:
```go
chargePointId := "cp0001"
csUrl = "ws://localhost:8887"
chargePoint := ocpp16.NewChargePoint(chargePointId, nil, nil)

// Set a handler for all callback functions
handler := &ChargePointHandler{}
chargePoint.SetChargePointCoreListener(handler)

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
request := ocpp16.NewBootNotificationRequest("model1", "vendor1")
```

You can then decide to send the message using a synchronous blocking call:
```go
// Synchronous call
confirmation, err := chargePoint.SendRequest(request)
if err != nil {
	log.Printf("error sending message: %v", err)
}
bootConf := confirmation.(*ocpp16.BootNotificationConfirmation)
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
callback := func(confirmation ocpp.Confirmation, e error) {
	bootConf := confirmation.(*ocpp16.BootNotificationConfirmation)
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
You can take a look at the full example inside `charge_point_sim.go`.
To run it, simply execute:
```bash
CLIENT_ID=chargePointSim CENTRAL_SYSTEM_URL=ws://<host>:8887 go run ./example/cp/charge_point_sim.go
```

You need to specify the hostname/IP of a running central station server, so the charge point can reach it.

#### Docker

A containerized version of the charge point example is available:
```bash
docker pull ldonini/ocpp1.6-charge-point:latest
docker run -e CLIENT_ID=chargePointSim -e CENTRAL_SYSTEM_URL=ws://<host>:8887 -it --rm --name charge-point ldonini/ocpp1.6-charge-point:latest
```

You need to specify the host, on which the central system is running, in order for the charge point to connect to it.

You can also build the docker image from source, using:
```sh
docker-compose up charge_point
```

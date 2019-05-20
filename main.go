package main

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/lorenzodonini/go-ocpp/ws"
	"log"
	"os"
	"strconv"
)

func runCentralSystem(args []string) {
	wsServer := ws.NewServer()
	centralSystem := ocpp.NewCentralSystem(wsServer, v16.CoreProfile.Profile)
	centralSystem.SetNewChargePointHandler(func(chargePointId string) {
		log.Printf("New charge point %v connected", chargePointId)
	})
	centralSystem.SetCallHandler(func(chargePointId string, call *ocpp.Call) {
		log.Printf("Call received from charge point %v: %v", chargePointId, call)
	})
	centralSystem.SetCallResultHandler(func(chargePointId string, callResult *ocpp.CallResult) {
		log.Printf("Call result received from charge point %v: %v", chargePointId, callResult)
	})
	centralSystem.SetCallErrorHandler(func(chargePointId string, callError *ocpp.CallError) {
		log.Printf("Call error received from charge point %v: %v", chargePointId, callError)
	})
	log.Print("Starting central system...")
	var listenPort int
	if len(args) > 1 {
		port, err := strconv.Atoi(args[1])
		if err != nil {
			listenPort = port
		}
	} else {
		listenPort = 8887
	}
	centralSystem.Start(listenPort, "/{ws}")
	log.Print("Stopped central system")
}

func runChargePoint(args []string) {
	wsClient := ws.NewClient()
	if len(args) != 3 {
		log.Print("Invalid client: chargePointId and centralSystemUrl are required")
		log.Print("Usage:\n\tocpp server [listenPort]\n\tocpp client id")
		return
	}
	id := args[1]
	csUrl := args[2]
	chargePoint := ocpp.NewChargePoint(id, wsClient, v16.CoreProfile.Profile)
	chargePoint.SetCallHandler(func(call *ocpp.Call) {
		log.Printf("Call received from central system: %v", call)
	})
	chargePoint.SetCallResultHandler(func(callResult *ocpp.CallResult) {
		log.Printf("Call result received from central system: %v", callResult)
	})
	chargePoint.SetCallErrorHandler(func(callError *ocpp.CallError) {
		log.Printf("Call error received from central system: %v", callError)
	})
	err := chargePoint.Start(csUrl)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Print("Usage:\n\tocpp server [listenPort]\n\tocpp client chargePointId centralSystemUrl")
	}
	mode := args[0]
	if mode == "client" {
		runChargePoint(args)
	} else if mode == "server" {
		runCentralSystem(args)
	} else {
		log.Print("Invalid mode")
		log.Print("Usage:\n\tocpp server [listenPort]\n\tocpp client [id]")
	}
}

package main

import (
	ocpp16 "github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"log"
	"os"
	"strconv"
)

type CentralSystemListener struct {
	chargePoints map[string]string
}

func (csl CentralSystemListener) OnAuthorize(chargePointId string, request *ocpp16.AuthorizeRequest) (confirmation *ocpp16.AuthorizeConfirmation, err error) {
	log.Printf("Received Authorize request from %v", chargePointId)
	return nil, nil
}

func (csl CentralSystemListener) OnBootNotification(chargePointId string, request *ocpp16.BootNotificationRequest) (confirmation *ocpp16.BootNotificationConfirmation, err error) {
	log.Printf("Received Boot Notification request from %v", chargePointId)
	return nil, nil
}

func runCentralSystem(args []string) {
	centralSystem := ocpp16.NewCentralSystem(nil, nil)
	listener := CentralSystemListener{chargePoints: map[string]string{}}
	centralSystem.SetNewChargePointHandler(func(chargePointId string) {
		log.Printf("New charge point %v connected", chargePointId)
	})
	centralSystem.SetCentralSystemCoreListener(listener)
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

type ChargePointListener struct {
}

func (cpl ChargePointListener) OnChangeAvailability(request *ocpp16.ChangeAvailabilityRequest) (confirmation *ocpp16.ChangeAvailabilityConfirmation, err error) {
	log.Printf("Received change availability request from central system")
	return nil, nil
}

func runChargePoint(args []string) {
	if len(args) != 3 {
		log.Print("Invalid client: chargePointId and centralSystemUrl are required")
		log.Print("Usage:\n\tocpp server [listenPort]\n\tocpp client id")
		return
	}
	id := args[1]
	csUrl := args[2]
	chargePoint := ocpp16.NewChargePoint(id, nil, nil)
	listener := ChargePointListener{}
	chargePoint.SetChargePointCoreListener(listener)
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

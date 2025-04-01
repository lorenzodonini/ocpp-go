package main

import (
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/sirupsen/logrus"
)

var csms ocpp2.CSMS
var server *ws.Server
var handler = &CsmsHandler{}

func main() {
	log.Info("Starting CSMS")
	csms.Start(7777, "/{charging_station_id}")
	log.Info("Stopped CSMS")
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	// Set this to DebugLevel if you want to retrieve verbose logs from the ocppj and websocket layers
	log.SetLevel(logrus.InfoLevel)

	server = ws.NewServer()
	server.SetBasicAuthHandler(handler.Authenticate)

	csms = ocpp2.NewCSMS(nil, server)
	csms.SetNewChargingStationHandler(handler.OnConnect)
	csms.SetChargingStationDisconnectedHandler(handler.OnDisconnect)
}

package main

import (
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/sirupsen/logrus"
)

const (
	listenPort = 7777
)

var log *logrus.Logger
var csms ocpp2.CSMS

func main() {
	// Run CSMS
	log.Infof("CSMS started on port %v", listenPort)
	csms.Start(listenPort, "/{charging_station_id}")
	log.Info("Stopped CSMS")
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	// Set this to DebugLevel if you want to retrieve verbose logs from the ocppj and websocket layers
	log.SetLevel(logrus.InfoLevel)

	csms = ocpp2.NewCSMS(nil, nil)
	csms.SetNewChargingStationHandler(func(cs ocpp2.ChargingStationConnection) {
		log.WithField("client", cs.ID()).Info("Charging Station connected")
	})
	csms.SetChargingStationDisconnectedHandler(func(cs ocpp2.ChargingStationConnection) {
		log.WithField("client", cs.ID()).Info("Charging Station disconnected")
	})
}

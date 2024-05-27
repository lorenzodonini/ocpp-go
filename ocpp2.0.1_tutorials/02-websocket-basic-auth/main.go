package main

import (
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/sirupsen/logrus"
)

const (
	listenPort = 7777
)

var log *logrus.Logger
var server *ws.Server
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
	ws.SetLogger(log.WithField("logger", "websocket"))
	ocppj.SetLogger(log.WithField("logger", "ocppj"))

	server = ws.NewServer()
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
}

func authenticate(user string, pass string) bool {
	return user == "cs001" && pass == "s3cr3t"
}

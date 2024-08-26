package main

import (
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/sirupsen/logrus"
)

const (
	csUrl = "ws://localhost:7777"
)

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

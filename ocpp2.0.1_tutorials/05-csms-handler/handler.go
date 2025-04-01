package main

import (
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type CsmsHandler struct {
}

func (csms *CsmsHandler) Authenticate(username, password string) bool {
	ok := username == "cs0001" && password == "s3cr3t"
	log.Info("Authenticated ", username, "? ", ok)
	return ok
}

func (handler *CsmsHandler) OnConnect(cs ocpp2.ChargingStationConnection) {
	log.WithField("client", cs.ID()).Info("Charging Station connected")
}

func (handler *CsmsHandler) OnDisconnect(cs ocpp2.ChargingStationConnection) {
	log.WithField("client", cs.ID()).Info("Charging Station disconnected")
}

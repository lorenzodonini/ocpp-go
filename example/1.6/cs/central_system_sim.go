package main

import (
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	defaultListenPort          = 8887
	defaultHeartbeatInterval   = 600
	envVarServerPort           = "SERVER_LISTEN_PORT"
	envVarTls                  = "TLS_ENABLED"
	envVarServerCertificate    = "SERVER_CERTIFICATE_PATH"
	envVarServerCertificateKey = "SERVER_CERTIFICATE_KEY_PATH"
)

func setupCentralSystem() ocpp16.CentralSystem {
	return ocpp16.NewCentralSystem(nil, nil)
}

func setupTlsCentralSystem() ocpp16.CentralSystem {
	certificate, ok := os.LookupEnv(envVarServerCertificate)
	if !ok {
		log.Fatalf("no required %v found", envVarServerCertificate)
	}
	key, ok := os.LookupEnv(envVarServerCertificateKey)
	if !ok {
		log.Fatalf("no required %v found", envVarServerCertificateKey)
	}
	server := ws.NewTLSServer(certificate, key)
	return ocpp16.NewCentralSystem(nil, server)
}

// Start function
func main() {
	// Load config from ENV
	var listenPort = defaultListenPort
	port, _ := os.LookupEnv(envVarServerPort)
	if p, err := strconv.Atoi(port); err == nil {
		listenPort = p
	} else {
		log.Printf("no valid %v environment variable found, using default port", envVarServerPort)
	}
	// Check if TLS enabled
	t, _ := os.LookupEnv(envVarTls)
	tlsEnabled, _ := strconv.ParseBool(t)
	// Prepare OCPP 1.6 central system
	var centralSystem ocpp16.CentralSystem
	if tlsEnabled {
		centralSystem = setupTlsCentralSystem()
	} else {
		centralSystem = setupCentralSystem()
	}
	// Support callbacks for all OCPP 1.6 profiles
	handler := &CentralSystemHandler{chargePoints: map[string]*ChargePointState{}}
	centralSystem.SetCoreHandler(handler)
	centralSystem.SetLocalAuthListHandler(handler)
	centralSystem.SetFirmwareManagementHandler(handler)
	centralSystem.SetReservationHandler(handler)
	centralSystem.SetRemoteTriggerHandler(handler)
	centralSystem.SetSmartChargingHandler(handler)
	// Add handlers for dis/connection of charge points
	centralSystem.SetNewChargePointHandler(func(chargePointId string) {
		handler.chargePoints[chargePointId] = &ChargePointState{connectors: map[int]*ConnectorInfo{}, transactions: map[int]*TransactionInfo{}}
		log.WithField("client", chargePointId).Info("new charge point connected")
	})
	centralSystem.SetChargePointDisconnectedHandler(func(chargePointId string) {
		log.WithField("client", chargePointId).Info("charge point disconnected")
		delete(handler.chargePoints, chargePointId)
	})
	// Run central system
	log.Infof("starting central system on port %v", listenPort)
	centralSystem.Start(listenPort, "/{ws}")
	log.Info("stopped central system")
}

func init() {
	log.SetLevel(log.InfoLevel)
}

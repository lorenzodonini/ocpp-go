package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

const (
	defaultListenPort          = 8887
	defaultHeartbeatInterval   = 600
	envVarServerPort           = "SERVER_LISTEN_PORT"
	envVarTls                  = "TLS_ENABLED"
	envVarCaCertificate        = "CA_CERTIFICATE_PATH"
	envVarServerCertificate    = "SERVER_CERTIFICATE_PATH"
	envVarServerCertificateKey = "SERVER_CERTIFICATE_KEY_PATH"
)

var log *logrus.Logger
var csms ocpp2.CSMS

func setupCentralSystem() ocpp2.CSMS {
	return ocpp2.NewCSMS(nil, nil)
}

func setupTlsCentralSystem() ocpp2.CSMS {
	var certPool *x509.CertPool
	// Load CA certificates
	caCertificate, ok := os.LookupEnv(envVarCaCertificate)
	if !ok {
		log.Infof("no %v found, using system CA pool", envVarCaCertificate)
		systemPool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("couldn't get system CA pool: %v", err)
		}
		certPool = systemPool
	} else {
		certPool = x509.NewCertPool()
		data, err := os.ReadFile(caCertificate)
		if err != nil {
			log.Fatalf("couldn't read CA certificate from %v: %v", caCertificate, err)
		}
		ok = certPool.AppendCertsFromPEM(data)
		if !ok {
			log.Fatalf("couldn't read CA certificate from %v", caCertificate)
		}
	}
	certificate, ok := os.LookupEnv(envVarServerCertificate)
	if !ok {
		log.Fatalf("no required %v found", envVarServerCertificate)
	}
	key, ok := os.LookupEnv(envVarServerCertificateKey)
	if !ok {
		log.Fatalf("no required %v found", envVarServerCertificateKey)
	}
	server := ws.NewServer(ws.WithServerTLSConfig(certificate, key, &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  certPool,
	}))
	return ocpp2.NewCSMS(nil, server)
}

// Run for every connected Charging Station, to simulate some functionality
func exampleRoutine(chargingStationID string, handler *CSMSHandler) {
	// Wait for some time
	time.Sleep(2 * time.Second)
	// Reserve a connector
	reservationID := 42
	clientIDTokenType := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	clientIdTag := "l33t"
	connectorID := 1
	expiryDate := types.NewDateTime(time.Now().Add(1 * time.Hour))
	cb1 := func(confirmation *reservation.ReserveNowResponse, err error) {
		if err != nil {
			logDefault(chargingStationID, reservation.ReserveNowFeatureName).Errorf("error on request: %v", err)
		} else if confirmation.Status == reservation.ReserveNowStatusAccepted {
			logDefault(chargingStationID, confirmation.GetFeatureName()).Infof("connector %v reserved for client %v until %v (reservation ID %d)", connectorID, clientIdTag, expiryDate.FormatTimestamp(), reservationID)
		} else {
			logDefault(chargingStationID, confirmation.GetFeatureName()).Infof("couldn't reserve connector %v: %v", connectorID, confirmation.Status)
		}
	}
	e := csms.ReserveNow(chargingStationID, cb1, reservationID, expiryDate, clientIDTokenType)
	if e != nil {
		logDefault(chargingStationID, reservation.ReserveNowFeatureName).Errorf("couldn't send message: %v", e)
		return
	}
	// Wait for some time
	time.Sleep(1 * time.Second)
	// Cancel the reservation
	cb2 := func(confirmation *reservation.CancelReservationResponse, err error) {
		if err != nil {
			logDefault(chargingStationID, reservation.CancelReservationFeatureName).Errorf("error on request: %v", err)
		} else if confirmation.Status == reservation.CancelReservationStatusAccepted {
			logDefault(chargingStationID, confirmation.GetFeatureName()).Infof("reservation %v canceled successfully", reservationID)
		} else {
			logDefault(chargingStationID, confirmation.GetFeatureName()).Infof("couldn't cancel reservation %v", reservationID)
		}
	}
	e = csms.CancelReservation(chargingStationID, cb2, reservationID)
	if e != nil {
		logDefault(chargingStationID, reservation.ReserveNowFeatureName).Errorf("couldn't send message: %v", e)
		return
	}
	// Wait for some time
	time.Sleep(5 * time.Second)
	// Get current local list version
	cb3 := func(confirmation *localauth.GetLocalListVersionResponse, err error) {
		if err != nil {
			logDefault(chargingStationID, localauth.GetLocalListVersionFeatureName).Errorf("error on request: %v", err)
		} else {
			logDefault(chargingStationID, confirmation.GetFeatureName()).Infof("current local list version: %v", confirmation.VersionNumber)
		}
	}
	e = csms.GetLocalListVersion(chargingStationID, cb3)
	if e != nil {
		logDefault(chargingStationID, localauth.GetLocalListVersionFeatureName).Errorf("couldn't send message: %v", e)
		return
	}
	// Wait for some time
	time.Sleep(5 * time.Second)
	setVariableData := []provisioning.SetVariableData{
		{
			AttributeType:  types.AttributeTarget,
			AttributeValue: "10",
			Component:      types.Component{Name: "OCPPCommCtrlr"},
			Variable:       types.Variable{Name: "HeartbeatInterval"},
		},
		{
			AttributeType:  types.AttributeTarget,
			AttributeValue: "true",
			Component:      types.Component{Name: "AuthCtrlr"},
			Variable:       types.Variable{Name: "Enabled"},
		},
	}
	// Change meter sampling values time
	cb4 := func(response *provisioning.SetVariablesResponse, err error) {
		if err != nil {
			logDefault(chargingStationID, provisioning.SetVariablesFeatureName).Errorf("error on request: %v", err)
			return
		}
		for _, r := range response.SetVariableResult {
			if r.AttributeStatus == provisioning.SetVariableStatusNotSupported {
				logDefault(chargingStationID, response.GetFeatureName()).Warnf("couldn't update variable %v for component %v: unsupported", r.Variable.Name, r.Component.Name)
			} else if r.AttributeStatus == provisioning.SetVariableStatusUnknownComponent {
				logDefault(chargingStationID, response.GetFeatureName()).Warnf("couldn't update variable for unknown component %v", r.Component.Name)
			} else if r.AttributeStatus == provisioning.SetVariableStatusUnknownVariable {
				logDefault(chargingStationID, response.GetFeatureName()).Warnf("couldn't update unknown variable %v for component %v", r.Variable.Name, r.Component.Name)
			} else if r.AttributeStatus == provisioning.SetVariableStatusRejected {
				logDefault(chargingStationID, response.GetFeatureName()).Warnf("couldn't update variable %v for key: %v", r.Variable.Name, r.Component.Name)
			} else {
				logDefault(chargingStationID, response.GetFeatureName()).Infof("updated variable %v for component %v", r.Variable.Name, r.Component.Name)
			}
		}
	}
	e = csms.SetVariables(chargingStationID, cb4, setVariableData)
	if e != nil {
		logDefault(chargingStationID, localauth.GetLocalListVersionFeatureName).Errorf("couldn't send message: %v", e)
		return
	}

	// Wait for some time
	time.Sleep(5 * time.Second)
	// Trigger a heartbeat message
	cb5 := func(response *remotecontrol.TriggerMessageResponse, err error) {
		if err != nil {
			logDefault(chargingStationID, remotecontrol.TriggerMessageFeatureName).Errorf("error on request: %v", err)
		} else if response.Status == remotecontrol.TriggerMessageStatusAccepted {
			logDefault(chargingStationID, response.GetFeatureName()).Infof("%v triggered successfully", availability.HeartbeatFeatureName)
		} else if response.Status == remotecontrol.TriggerMessageStatusRejected {
			logDefault(chargingStationID, response.GetFeatureName()).Infof("%v trigger was rejected", availability.HeartbeatFeatureName)
		}
	}
	e = csms.TriggerMessage(chargingStationID, cb5, remotecontrol.MessageTriggerHeartbeat)
	if e != nil {
		logDefault(chargingStationID, remotecontrol.TriggerMessageFeatureName).Errorf("couldn't send message: %v", e)
		return
	}

	// Wait for some time
	time.Sleep(5 * time.Second)
	// Trigger a diagnostics status notification
	cb6 := func(response *remotecontrol.TriggerMessageResponse, err error) {
		if err != nil {
			logDefault(chargingStationID, remotecontrol.TriggerMessageFeatureName).Errorf("error on request: %v", err)
		} else if response.Status == remotecontrol.TriggerMessageStatusAccepted {
			logDefault(chargingStationID, response.GetFeatureName()).Infof("%v triggered successfully", diagnostics.LogStatusNotificationFeatureName)
		} else if response.Status == remotecontrol.TriggerMessageStatusRejected {
			logDefault(chargingStationID, response.GetFeatureName()).Infof("%v trigger was rejected", diagnostics.LogStatusNotificationFeatureName)
		}
	}
	e = csms.TriggerMessage(chargingStationID, cb6, remotecontrol.MessageTriggerLogStatusNotification)
	if e != nil {
		logDefault(chargingStationID, remotecontrol.TriggerMessageFeatureName).Errorf("couldn't send message: %v", e)
		return
	}

	// Wait for some time
	time.Sleep(5 * time.Second)
	// Set a custom display message
	cb7 := func(response *display.SetDisplayMessageResponse, err error) {
		if err != nil {
			logDefault(chargingStationID, display.SetDisplayMessageFeatureName).Errorf("error on request: %v", err)
		} else if response.Status == display.DisplayMessageStatusAccepted {
			logDefault(chargingStationID, response.GetFeatureName()).Info("display message set successfully")
		} else {
			logDefault(chargingStationID, response.GetFeatureName()).Errorf("failed to set display message: %v", response.Status)
		}
	}
	var currentTx int
	for txID := range handler.chargingStations[chargingStationID].transactions {
		currentTx = txID
		break
	}
	e = csms.SetDisplayMessage(chargingStationID, cb7, display.MessageInfo{
		ID:            42,
		Priority:      display.MessagePriorityInFront,
		State:         display.MessageStateCharging,
		TransactionID: fmt.Sprintf("%d", currentTx),
		Message: types.MessageContent{
			Format:   types.MessageFormatUTF8,
			Language: "en-US",
			Content:  "Hello world!",
		},
	})
	if e != nil {
		logDefault(chargingStationID, display.SetDisplayMessageFeatureName).Errorf("couldn't send message: %v", e)
		return
	}
	// Finish simulation
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
	if tlsEnabled {
		csms = setupTlsCentralSystem()
	} else {
		csms = setupCentralSystem()
	}
	// Support callbacks for all OCPP 2.0.1 profiles
	handler := &CSMSHandler{chargingStations: map[string]*ChargingStationState{}}
	csms.SetAuthorizationHandler(handler)
	csms.SetAvailabilityHandler(handler)
	csms.SetDiagnosticsHandler(handler)
	csms.SetFirmwareHandler(handler)
	csms.SetLocalAuthListHandler(handler)
	csms.SetMeterHandler(handler)
	csms.SetProvisioningHandler(handler)
	csms.SetRemoteControlHandler(handler)
	csms.SetReservationHandler(handler)
	csms.SetTariffCostHandler(handler)
	csms.SetTransactionsHandler(handler)
	// Add handlers for dis/connection of charging stations
	csms.SetNewChargingStationHandler(func(chargingStation ocpp2.ChargingStationConnection) {
		handler.chargingStations[chargingStation.ID()] = &ChargingStationState{connectors: map[int]*ConnectorInfo{}, transactions: map[int]*TransactionInfo{}}
		log.WithField("client", chargingStation.ID()).Info("new charging station connected")
		go exampleRoutine(chargingStation.ID(), handler)
	})
	csms.SetChargingStationDisconnectedHandler(func(chargingStation ocpp2.ChargingStationConnection) {
		log.WithField("client", chargingStation.ID()).Info("charging station disconnected")
		delete(handler.chargingStations, chargingStation.ID())
	})
	ocppj.SetLogger(log)
	// Run CSMS
	log.Infof("starting CSMS on port %v", listenPort)
	csms.Start(listenPort, "/{ws}")
	log.Info("stopped CSMS")
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetLevel(logrus.InfoLevel)
}

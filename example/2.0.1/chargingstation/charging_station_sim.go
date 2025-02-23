package main

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"strconv"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/transactions"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/sirupsen/logrus"

	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

const (
	envVarClientID             = "CLIENT_ID"
	envVarCSMSUrl              = "CSMS_URL"
	envVarTls                  = "TLS_ENABLED"
	envVarCACertificate        = "CA_CERTIFICATE_PATH"
	envVarClientCertificate    = "CLIENT_CERTIFICATE_PATH"
	envVarClientCertificateKey = "CLIENT_CERTIFICATE_KEY_PATH"
)

var log *logrus.Logger

func setupChargingStation(chargingStationID string) ocpp2.ChargingStation {
	return ocpp2.NewChargingStation(chargingStationID, nil, nil)
}

func setupTlsChargingStation(chargingStationID string) ocpp2.ChargingStation {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	// Load CA cert
	caPath, ok := os.LookupEnv(envVarCACertificate)
	if ok {
		caCert, err := os.ReadFile(caPath)
		if err != nil {
			log.Warn(err)
		} else if !certPool.AppendCertsFromPEM(caCert) {
			log.Info("no ca.cert file found, will use system CA certificates")
		}
	} else {
		log.Info("no ca.cert file found, will use system CA certificates")
	}
	// Load client certificate
	clientCertPath, ok1 := os.LookupEnv(envVarClientCertificate)
	clientKeyPath, ok2 := os.LookupEnv(envVarClientCertificateKey)
	var clientCertificates []tls.Certificate
	if ok1 && ok2 {
		certificate, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err == nil {
			clientCertificates = []tls.Certificate{certificate}
		} else {
			log.Infof("couldn't load client TLS certificate: %v", err)
		}
	}
	// Create client with TLS config
	client := ws.NewClient(ws.WithClientTLSConfig(&tls.Config{
		RootCAs:      certPool,
		Certificates: clientCertificates,
	}))
	return ocpp2.NewChargingStation(chargingStationID, nil, client)
}

// exampleRoutine simulates a charging station flow, where a dummy transaction is started.
// The simulation runs for about 5 minutes.
func exampleRoutine(chargingStation ocpp2.ChargingStation, stateHandler *ChargingStationHandler) {
	dummyClientIdToken := types.IdToken{
		IdToken: "12345",
		Type:    types.IdTokenTypeKeyCode,
	}
	// Boot
	bootResp, err := chargingStation.BootNotification(provisioning.BootReasonPowerUp, "model1", "vendor1")
	checkError(err)
	logDefault(bootResp.GetFeatureName()).Infof("status: %v, interval: %v, current time: %v", bootResp.Status, bootResp.Interval, bootResp.CurrentTime.String())
	// Notify EVSE status
	for eID, e := range stateHandler.evse {
		updateOperationalStatus(stateHandler, eID, availability.OperationalStatusOperative)
		// Notify connector status
		for cID := range e.connectors {
			updateConnectorStatus(stateHandler, eID, cID, availability.ConnectorStatusAvailable)
		}
	}
	// Wait for some time ...
	time.Sleep(5 * time.Second)
	// Simulate charging for connector 1
	// EV is plugged in
	evseID := 1
	evse := stateHandler.evse[evseID]
	chargingConnector := 0
	updateConnectorStatus(stateHandler, evseID, chargingConnector, availability.ConnectorStatusOccupied)
	// Start transaction
	tx := transactions.Transaction{
		TransactionID: pseudoUUID(),
		ChargingState: transactions.ChargingStateEVConnected,
	}
	evseReq := types.EVSE{ID: evseID, ConnectorID: &chargingConnector}
	txEventResp, err := chargingStation.TransactionEvent(transactions.TransactionEventStarted, types.Now(), transactions.TriggerReasonCablePluggedIn, evse.nextSequence(), tx, func(request *transactions.TransactionEventRequest) {
		request.Evse = &evseReq
	})
	checkError(err)
	logDefault(txEventResp.GetFeatureName()).Infof("transaction %v started", tx.TransactionID)
	stateHandler.evse[evseID].currentTransaction = tx.TransactionID
	// Authorize
	authResp, err := chargingStation.Authorize(dummyClientIdToken.IdToken, types.IdTokenTypeKeyCode)
	checkError(err)
	logDefault(authResp.GetFeatureName()).Infof("status: %v %v", authResp.IdTokenInfo.Status, getExpiryDate(&authResp.IdTokenInfo))
	// Update transaction with auth info
	txEventResp, err = chargingStation.TransactionEvent(transactions.TransactionEventUpdated, types.Now(), transactions.TriggerReasonAuthorized, evse.nextSequence(), tx, func(request *transactions.TransactionEventRequest) {
		request.Evse = &evseReq
		request.IDToken = &dummyClientIdToken
	})
	checkError(err)
	logDefault(txEventResp.GetFeatureName()).Infof("transaction %v updated", tx.TransactionID)
	// Update transaction after energy offering starts
	txEventResp, err = chargingStation.TransactionEvent(transactions.TransactionEventUpdated, types.Now(), transactions.TriggerReasonChargingStateChanged, evse.nextSequence(), tx, func(request *transactions.TransactionEventRequest) {
		request.Evse = &evseReq
		request.IDToken = &dummyClientIdToken
	})
	checkError(err)
	logDefault(txEventResp.GetFeatureName()).Infof("transaction %v updated", tx.TransactionID)
	// Periodically send meter values
	var sampleInterval time.Duration = 5
	//sampleInterval, ok := stateHandler.configuration.getInt(MeterValueSampleInterval)
	//if !ok {
	//	sampleInterval = 5
	//}
	var sampledValue types.SampledValue
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * sampleInterval)
		stateHandler.meterValue += 10
		sampledValue = types.SampledValue{
			Value:     stateHandler.meterValue,
			Context:   types.ReadingContextSamplePeriodic,
			Measurand: types.MeasurandEnergyActiveExportRegister,
			Phase:     types.PhaseL3,
			Location:  types.LocationOutlet,
			UnitOfMeasure: &types.UnitOfMeasure{
				Unit: "kWh",
			},
		}
		meterValue := types.MeterValue{
			Timestamp:    types.DateTime{Time: time.Now()},
			SampledValue: []types.SampledValue{sampledValue},
		}
		// Send meter values
		txEventResp, err = chargingStation.TransactionEvent(transactions.TransactionEventUpdated, types.Now(), transactions.TriggerReasonMeterValuePeriodic, evse.nextSequence(), tx, func(request *transactions.TransactionEventRequest) {
			request.MeterValue = []types.MeterValue{meterValue}
			request.IDToken = &dummyClientIdToken
		})
		checkError(err)
		logDefault(txEventResp.GetFeatureName()).Infof("transaction %v updated with periodic meter values", tx.TransactionID)
		// Increase meter value
		stateHandler.meterValue += 2
	}
	// Stop charging for connector 1
	updateConnectorStatus(stateHandler, evseID, chargingConnector, availability.ConnectorStatusAvailable)
	// Send transaction end data
	sampledValue.Context = types.ReadingContextTransactionEnd
	sampledValue.Value = stateHandler.meterValue
	tx.StoppedReason = transactions.ReasonEVDisconnected
	txEventResp, err = chargingStation.TransactionEvent(transactions.TransactionEventEnded, types.Now(), transactions.TriggerReasonEVCommunicationLost, evse.nextSequence(), tx, func(request *transactions.TransactionEventRequest) {
		request.Evse = &evseReq
		request.IDToken = &dummyClientIdToken
		request.MeterValue = []types.MeterValue{}
	})
	checkError(err)
	logDefault(txEventResp.GetFeatureName()).Infof("transaction %v stopped", tx.TransactionID)
	// Wait for some time ...
	time.Sleep(5 * time.Minute)
	// End simulation
}

// Start function
func main() {
	// Load config
	id, ok := os.LookupEnv(envVarClientID)
	if !ok {
		log.Printf("no %v environment variable found, exiting...", envVarClientID)
		return
	}
	csmsUrl, ok := os.LookupEnv(envVarCSMSUrl)
	if !ok {
		log.Printf("no %v environment variable found, exiting...", envVarCSMSUrl)
		return
	}
	// Check if TLS enabled
	t, _ := os.LookupEnv(envVarTls)
	tlsEnabled, _ := strconv.ParseBool(t)
	// Prepare OCPP 2.0.1 charging station (chargingStation variable is defined in handler.go)
	if tlsEnabled {
		chargingStation = setupTlsChargingStation(id)
	} else {
		chargingStation = setupChargingStation(id)
	}
	// Setup some basic state management
	evse := EVSEInfo{
		availability:       availability.OperationalStatusOperative,
		currentTransaction: "",
		currentReservation: 0,
		connectors: map[int]ConnectorInfo{
			0: {
				status:       availability.ConnectorStatusAvailable,
				availability: availability.OperationalStatusOperative,
				typ:          reservation.ConnectorTypeCType2,
			},
		},
		seqNo: 0,
	}
	handler := &ChargingStationHandler{
		model:                "model1",
		vendor:               "vendor1",
		availability:         availability.OperationalStatusOperative,
		evse:                 map[int]*EVSEInfo{1: &evse},
		localAuthList:        []localauth.AuthorizationData{},
		localAuthListVersion: 0,
		monitoringLevel:      0,
		meterValue:           0,
	}
	// Support callbacks for all OCPP 2.0.1 profiles
	chargingStation.SetAvailabilityHandler(handler)
	chargingStation.SetAuthorizationHandler(handler)
	chargingStation.SetDataHandler(handler)
	chargingStation.SetDiagnosticsHandler(handler)
	chargingStation.SetDisplayHandler(handler)
	chargingStation.SetFirmwareHandler(handler)
	chargingStation.SetISO15118Handler(handler)
	chargingStation.SetLocalAuthListHandler(handler)
	chargingStation.SetProvisioningHandler(handler)
	chargingStation.SetRemoteControlHandler(handler)
	chargingStation.SetReservationHandler(handler)
	chargingStation.SetSmartChargingHandler(handler)
	chargingStation.SetTariffCostHandler(handler)
	chargingStation.SetTransactionsHandler(handler)
	ocppj.SetLogger(log)
	// Connects to central system
	err := chargingStation.Start(csmsUrl)
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("connected to CSMS at %v", csmsUrl)
		exampleRoutine(chargingStation, handler)
		// Disconnect
		chargingStation.Stop()
		log.Infof("disconnected from CSMS")
	}
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetLevel(logrus.InfoLevel)
}

// Utility functions
func logDefault(feature string) *logrus.Entry {
	return log.WithField("message", feature)
}

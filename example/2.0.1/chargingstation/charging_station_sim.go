package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

const (
	envVarClientID             = "CLIENT_ID"
	envVarCentralSystemUrl     = "CENTRAL_SYSTEM_URL"
	envVarTls                  = "TLS_ENABLED"
	envVarCACertificate        = "CA_CERTIFICATE_PATH"
	envVarClientCertificate    = "CLIENT_CERTIFICATE_PATH"
	envVarClientCertificateKey = "CLIENT_CERTIFICATE_KEY_PATH"
)

var log *logrus.Logger

func setupChargePoint(chargePointID string) ocpp2.ChargingStation {
	return ocpp2.NewChargingStation(chargePointID, nil, nil)
}

func setupTlsChargePoint(chargePointID string) ocpp2.ChargingStation {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	// Load CA cert
	caPath, ok := os.LookupEnv(envVarCACertificate)
	if ok {
		caCert, err := ioutil.ReadFile(caPath)
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
	client := ws.NewTLSClient(&tls.Config{
		RootCAs:      certPool,
		Certificates: clientCertificates,
	})
	return ocpp2.NewChargingStation(chargePointID, nil, client)
}

// exampleRoutine simulates a charge point flow, where
func exampleRoutine(chargingStation ocpp2.ChargingStation, stateHandler *ChargingStationHandler) {
	dummyClientIdTag := "12345"
	chargingConnector := 1
	// Boot
	bootConf, err := chargingStation.BootNotification(provisioning.BootReasonPowerUp, "model1", "vendor1")
	checkError(err)
	logDefault(bootConf.GetFeatureName()).Infof("status: %v, interval: %v, current time: %v", bootConf.Status, bootConf.Interval, bootConf.CurrentTime.String())
	// Notify connector status
	updateStatus(stateHandler, 0, core.ChargePointStatusAvailable)
	// Wait for some time ...
	time.Sleep(5 * time.Second)
	// Simulate charging for connector 1
	authResp, err := chargingStation.Authorize(dummyClientIdTag, types.IdTokenTypeKeyCode)
	checkError(err)
	logDefault(authResp.GetFeatureName()).Infof("status: %v %v", authResp.IdTokenInfo.Status, getExpiryDate(authResp.IdTokenInfo))
	// Update connector status
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusPreparing)
	// Start transaction
	startConf, err := chargingStation.StartTransaction(chargingConnector, dummyClientIdTag, stateHandler.meterValue, types.NewDateTime(time.Now()))
	checkError(err)
	logDefault(startConf.GetFeatureName()).Infof("status: %v, transaction %v %v", startConf.IdTagInfo.Status, startConf.TransactionId, getExpiryDate(startConf.IdTagInfo))
	stateHandler.connectors[chargingConnector].currentTransaction = startConf.TransactionId
	// Update connector status
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusCharging)
	// Periodically send meter values
	for i := 0; i < 5; i++ {
		sampleInterval, ok := stateHandler.configuration.getInt(MeterValueSampleInterval)
		if !ok {
			sampleInterval = 5
		}
		time.Sleep(time.Second * time.Duration(sampleInterval))
		stateHandler.meterValue += 10
		sampledValue := types.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit: types.UnitOfMeasureWh, Format: types.ValueFormatRaw, Measurand: types.MeasurandEnergyActiveExportRegister, Context: types.ReadingContextSamplePeriodic, Location: types.LocationOutlet}
		meterValue := types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{sampledValue}}
		meterConf, err := chargingStation.MeterValues(chargingConnector, []types.MeterValue{meterValue})
		checkError(err)
		logDefault(meterConf.GetFeatureName()).Infof("sent updated %v", sampledValue.Measurand)
	}
	stateHandler.meterValue += 2
	// Stop charging for connector 1
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusFinishing)
	stopConf, err := chargingStation.StopTransaction(stateHandler.meterValue, types.NewDateTime(time.Now()), startConf.TransactionId, func(request *core.StopTransactionRequest) {
		sampledValue := types.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit: types.UnitOfMeasureWh, Format: types.ValueFormatRaw, Measurand: types.MeasurandEnergyActiveExportRegister, Context: types.ReadingContextSamplePeriodic, Location: types.LocationOutlet}
		meterValue := types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{sampledValue}}
		request.TransactionData = []types.MeterValue{meterValue}
		request.Reason = core.ReasonEVDisconnected
	})
	checkError(err)
	logDefault(stopConf.GetFeatureName()).Infof("transaction %v stopped", startConf.TransactionId)
	// Update connector status
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusAvailable)
	// Wait for some time ...
	time.Sleep(5 * time.Minute)
}

// Start function
func main() {
	// Load config
	id, ok := os.LookupEnv(envVarClientID)
	if !ok {
		log.Printf("no %v environment variable found, exiting...", envVarClientID)
		return
	}
	csUrl, ok := os.LookupEnv(envVarCentralSystemUrl)
	if !ok {
		log.Printf("no %v environment variable found, exiting...", envVarCentralSystemUrl)
		return
	}
	// Check if TLS enabled
	t, _ := os.LookupEnv(envVarTls)
	tlsEnabled, _ := strconv.ParseBool(t)
	// Prepare OCPP 1.6 charge point (chargePoint variable is defined in handler.go)
	if tlsEnabled {
		chargingStation = setupTlsChargePoint(id)
	} else {
		chargingStation = setupChargePoint(id)
	}
	// Setup some basic state management
	connectors := map[int]*ConnectorInfo{
		1: {status: core.ChargePointStatusAvailable, availability: core.AvailabilityTypeOperative, currentTransaction: 0},
	}
	handler := &ChargingStationHandler{
		status:               core.ChargePointStatusAvailable,
		connectors:           connectors,
		configuration:        getDefaultConfig(),
		errorCode:            core.NoError,
		localAuthList:        []localauth.AuthorizationData{},
		localAuthListVersion: 0}
	// Support callbacks for all OCPP 1.6 profiles
	chargingStation.SetCoreHandler(handler)
	chargingStation.SetFirmwareManagementHandler(handler)
	chargingStation.SetLocalAuthListHandler(handler)
	chargingStation.SetReservationHandler(handler)
	chargingStation.SetRemoteTriggerHandler(handler)
	chargingStation.SetSmartChargingHandler(handler)
	ocppj.SetLogger(log)
	// Connects to central system
	err := chargingStation.Start(csUrl)
	if err != nil {
		log.Errorln(err)
	} else {
		log.Infof("connected to central system at %v", csUrl)
		exampleRoutine(chargingStation, handler)
		// Disconnect
		chargingStation.Stop()
		log.Infof("disconnected from central system")
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

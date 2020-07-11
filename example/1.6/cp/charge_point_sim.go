package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const (
	envVarClientID             = "CLIENT_ID"
	envVarCentralSystemUrl     = "CENTRAL_SYSTEM_URL"
	envVarTls                  = "TLS_ENABLED"
	envVarCACertificate        = "CA_CERTIFICATE_PATH"
	envVarClientCertificate    = "CLIENT_CERTIFICATE_PATH"
	envVarClientCertificateKey = "CLIENT_CERTIFICATE_KEY_PATH"
)

func setupChargePoint(chargePointID string) ocpp16.ChargePoint {
	return ocpp16.NewChargePoint(chargePointID, nil, nil)
}

func setupTlsChargePoint(chargePointID string) ocpp16.ChargePoint {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	// Load CA cert
	caPath, ok := os.LookupEnv(envVarCACertificate)
	if !ok {
		log.Fatalf("no required %v found", envVarCACertificate)
	}
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		log.Warn(err)
	} else if !certPool.AppendCertsFromPEM(caCert) {
		log.Warn("no ca.cert file found, will use system CA certificates")
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
	return ocpp16.NewChargePoint(chargePointID, nil, client)
}

// Used for scheduling requests that are not generated from the main thread.
// OCPP is a request-response protocol and doesn't support multiple simultaneous requests.
//
// By scheduling a request, it can later on be sent safely from the main thread.
func scheduleAsyncRequest(asyncRequest func()) {
	asyncRequestChan <- asyncRequest
}

// Wait is used for simulating idle time, while being able to process trigger requests.
//
// If trigger requests (or other asynchronous requests) are fired while waiting, the request will be sent.
// This mechanism prevents errors when attempting to send a request while another one is still pending.
func wait(d time.Duration) {
	t := time.NewTimer(d)
	for {
		select {
		case req, ok := <-asyncRequestChan:
			if !ok {
				return
			}
			req()
		case <-t.C:
			log.Debugf("finished simulated wait")
			return
		}
	}
}

// exampleRoutine simulates a charge point flow, where
func exampleRoutine(chargePoint ocpp16.ChargePoint, stateHandler *ChargePointHandler) {
	dummyClientIdTag := "12345"
	chargingConnector := 1
	// Boot
	bootConf, err := chargePoint.BootNotification("model1", "vendor1")
	checkError(err)
	logDefault(bootConf.GetFeatureName()).Infof("status: %v, interval: %v, current time: %v", bootConf.Status, bootConf.Interval, bootConf.CurrentTime.String())
	// Notify connector status
	updateStatus(stateHandler, 0, core.ChargePointStatusAvailable)
	// Wait for some time ...
	wait(5 * time.Second)
	// Simulate charging for connector 1
	authConf, err := chargePoint.Authorize(dummyClientIdTag)
	checkError(err)
	logDefault(authConf.GetFeatureName()).Infof("status: %v %v", authConf.IdTagInfo.Status, getExpiryDate(authConf.IdTagInfo))
	// Update connector status
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusPreparing)
	// Start transaction
	startConf, err := chargePoint.StartTransaction(chargingConnector, dummyClientIdTag, stateHandler.meterValue, types.NewDateTime(time.Now()))
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
		wait(time.Second * time.Duration(sampleInterval))
		stateHandler.meterValue += 10
		sampledValue := types.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit: types.UnitOfMeasureWh, Format: types.ValueFormatRaw, Measurand: types.MeasurandEnergyActiveExportRegister, Context: types.ReadingContextSamplePeriodic, Location: types.LocationOutlet}
		meterValue := types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{sampledValue}}
		meterConf, err := chargePoint.MeterValues(chargingConnector, []types.MeterValue{meterValue})
		checkError(err)
		logDefault(meterConf.GetFeatureName()).Infof("sent updated %v", sampledValue.Measurand)
	}
	stateHandler.meterValue += 2
	// Stop charging for connector 1
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusFinishing)
	stopConf, err := chargePoint.StopTransaction(stateHandler.meterValue, types.NewDateTime(time.Now()), startConf.TransactionId, func(request *core.StopTransactionRequest) {
		sampledValue := types.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit: types.UnitOfMeasureWh, Format: types.ValueFormatRaw, Measurand: types.MeasurandEnergyActiveExportRegister, Context: types.ReadingContextSamplePeriodic, Location: types.LocationOutlet}
		meterValue := types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{sampledValue}}
		request.TransactionData = []types.MeterValue{meterValue}
		request.Reason = core.ReasonEVDisconnected
	})
	checkError(err)
	logDefault(stopConf.GetFeatureName()).Infof("transaction %v stopped", startConf.TransactionId)
	// Update connector status
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusAvailable)
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
		chargePoint = setupTlsChargePoint(id)
	} else {
		chargePoint = setupChargePoint(id)
	}
	// Setup some basic state management
	connectors := map[int]*ConnectorInfo{
		1: {status: core.ChargePointStatusAvailable, availability: core.AvailabilityTypeOperative, currentTransaction: 0},
	}
	handler := &ChargePointHandler{
		status:               core.ChargePointStatusAvailable,
		connectors:           connectors,
		configuration:        getDefaultConfig(),
		errorCode:            core.NoError,
		localAuthList:        []localauth.AuthorizationData{},
		localAuthListVersion: 0}
	// Support callbacks for all OCPP 1.6 profiles
	chargePoint.SetCoreHandler(handler)
	chargePoint.SetFirmwareManagementHandler(handler)
	chargePoint.SetLocalAuthListHandler(handler)
	chargePoint.SetReservationHandler(handler)
	chargePoint.SetRemoteTriggerHandler(handler)
	chargePoint.SetSmartChargingHandler(handler)
	// Connects to central system
	err := chargePoint.Start(csUrl)
	if err != nil {
		log.Println(err)
	} else {
		log.Infof("connected to central system at %v", csUrl)
		// Setup channel for asynchronous requests (e.g. triggers)
		asyncRequestChan = make(chan func(), 5)
		exampleRoutine(chargePoint, handler)
		close(asyncRequestChan)
		// Disconnect
		chargePoint.Stop()
		log.Infof("disconnected from central system")
	}
}

func init() {
	log.SetLevel(log.InfoLevel)
}

// Utility functions
func logDefault(feature string) *log.Entry {
	return log.WithField("message", feature)
}

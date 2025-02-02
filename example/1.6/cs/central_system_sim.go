package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
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
	envVarMetricsEnabled       = "METRICS_ENABLED"
	envVarMetricsAddress       = "METRICS_ADDRESS"
)

var log *logrus.Logger
var centralSystem ocpp16.CentralSystem

func setupCentralSystem() ocpp16.CentralSystem {
	return ocpp16.NewCentralSystem(nil, nil)
}

func setupTlsCentralSystem() ocpp16.CentralSystem {
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
	return ocpp16.NewCentralSystem(nil, server)
}

// Run for every connected Charge Point, to simulate some functionality
func exampleRoutine(chargePointID string, handler *CentralSystemHandler) {
	// Wait for some time
	time.Sleep(2 * time.Second)
	// Reserve a connector
	reservationID := 42
	clientIdTag := "l33t"
	connectorID := 1
	expiryDate := types.NewDateTime(time.Now().Add(1 * time.Hour))
	cb1 := func(confirmation *reservation.ReserveNowConfirmation, err error) {
		if err != nil {
			logDefault(chargePointID, reservation.ReserveNowFeatureName).Errorf("error on request: %v", err)
		} else if confirmation.Status == reservation.ReservationStatusAccepted {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("connector %v reserved for client %v until %v (reservation ID %d)", connectorID, clientIdTag, expiryDate.FormatTimestamp(), reservationID)
		} else {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("couldn't reserve connector %v: %v", connectorID, confirmation.Status)
		}
	}
	e := centralSystem.ReserveNow(chargePointID, cb1, connectorID, expiryDate, clientIdTag, reservationID)
	if e != nil {
		logDefault(chargePointID, reservation.ReserveNowFeatureName).Errorf("couldn't send message: %v", e)
		return
	}
	// Wait for some time
	time.Sleep(1 * time.Second)
	// Cancel the reservation
	cb2 := func(confirmation *reservation.CancelReservationConfirmation, err error) {
		if err != nil {
			logDefault(chargePointID, reservation.CancelReservationFeatureName).Errorf("error on request: %v", err)
		} else if confirmation.Status == reservation.CancelReservationStatusAccepted {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("reservation %v canceled successfully", reservationID)
		} else {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("couldn't cancel reservation %v", reservationID)
		}
	}
	e = centralSystem.CancelReservation(chargePointID, cb2, reservationID)
	if e != nil {
		logDefault(chargePointID, reservation.ReserveNowFeatureName).Errorf("couldn't send message: %v", e)
		return
	}
	// Wait for some time
	time.Sleep(5 * time.Second)
	// Get current local list version
	cb3 := func(confirmation *localauth.GetLocalListVersionConfirmation, err error) {
		if err != nil {
			logDefault(chargePointID, localauth.GetLocalListVersionFeatureName).Errorf("error on request: %v", err)
		} else {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("current local list version: %v", confirmation.ListVersion)
		}
	}
	e = centralSystem.GetLocalListVersion(chargePointID, cb3)
	if e != nil {
		logDefault(chargePointID, localauth.GetLocalListVersionFeatureName).Errorf("couldn't send message: %v", e)
		return
	}
	// Wait for some time
	time.Sleep(5 * time.Second)
	configKey := "MeterValueSampleInterval"
	configValue := "10"
	// Change meter sampling values time
	cb4 := func(confirmation *core.ChangeConfigurationConfirmation, err error) {
		if err != nil {
			logDefault(chargePointID, core.ChangeConfigurationFeatureName).Errorf("error on request: %v", err)
		} else if confirmation.Status == core.ConfigurationStatusNotSupported {
			logDefault(chargePointID, confirmation.GetFeatureName()).Warnf("couldn't update configuration for unsupported key: %v", configKey)
		} else if confirmation.Status == core.ConfigurationStatusRejected {
			logDefault(chargePointID, confirmation.GetFeatureName()).Warnf("couldn't update configuration for readonly key: %v", configKey)
		} else {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("updated configuration for key %v to: %v", configKey, configValue)
		}
	}
	e = centralSystem.ChangeConfiguration(chargePointID, cb4, configKey, configValue)
	if e != nil {
		logDefault(chargePointID, localauth.GetLocalListVersionFeatureName).Errorf("couldn't send message: %v", e)
		return
	}

	// Wait for some time
	time.Sleep(5 * time.Second)
	// Trigger a heartbeat message
	cb5 := func(confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
		if err != nil {
			logDefault(chargePointID, remotetrigger.TriggerMessageFeatureName).Errorf("error on request: %v", err)
		} else if confirmation.Status == remotetrigger.TriggerMessageStatusAccepted {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("%v triggered successfully", core.HeartbeatFeatureName)
		} else if confirmation.Status == remotetrigger.TriggerMessageStatusRejected {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("%v trigger was rejected", core.HeartbeatFeatureName)
		}
	}
	e = centralSystem.TriggerMessage(chargePointID, cb5, core.HeartbeatFeatureName)
	if e != nil {
		logDefault(chargePointID, remotetrigger.TriggerMessageFeatureName).Errorf("couldn't send message: %v", e)
		return
	}

	// Wait for some time
	time.Sleep(5 * time.Second)
	// Trigger a diagnostics status notification
	cb6 := func(confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
		if err != nil {
			logDefault(chargePointID, remotetrigger.TriggerMessageFeatureName).Errorf("error on request: %v", err)
		} else if confirmation.Status == remotetrigger.TriggerMessageStatusAccepted {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("%v triggered successfully", firmware.GetDiagnosticsFeatureName)
		} else if confirmation.Status == remotetrigger.TriggerMessageStatusRejected {
			logDefault(chargePointID, confirmation.GetFeatureName()).Infof("%v trigger was rejected", firmware.GetDiagnosticsFeatureName)
		}
	}
	e = centralSystem.TriggerMessage(chargePointID, cb6, firmware.DiagnosticsStatusNotificationFeatureName)
	if e != nil {
		logDefault(chargePointID, remotetrigger.TriggerMessageFeatureName).Errorf("couldn't send message: %v", e)
		return
	}
}

// sets up OTLP metrics exporter
func setupMetrics(address string) error {
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	client, err := grpc.NewClient(address, grpcOpts...)

	if err != nil {
		return errors.Wrap(err, "failed to create gRPC connection to collector")
	}

	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(client))
	if err != nil {
		return errors.Wrap(err, "failed to create otlp metric exporter")
	}

	resource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("centralSystem-demo"),
			semconv.ServiceVersionKey.String("example"),
		),
		resource.WithFromEnv(),
		resource.WithContainer(),
		resource.WithOS(),
		resource.WithOSType(),
		resource.WithHost(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create resource")
	}

	meterProvider := metricsdk.NewMeterProvider(
		metricsdk.WithReader(
			metricsdk.NewPeriodicReader(
				exporter,
				metricsdk.WithInterval(10*time.Second),
			),
		),
		metricsdk.WithResource(resource),
	)

	otel.SetMeterProvider(meterProvider)
	return nil
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

	// Setup metrics if enabled
	if t, _ := os.LookupEnv(envVarMetricsEnabled); t == "true" {
		address, _ := os.LookupEnv(envVarMetricsAddress)
		if err := setupMetrics(address); err != nil {
			log.Error(err)
			return
		}
	}

	// Check if TLS enabled
	t, _ := os.LookupEnv(envVarTls)
	tlsEnabled, _ := strconv.ParseBool(t)
	// Prepare OCPP 1.6 central system
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

	// Add callbacks for OCPP 1.6 security profiles
	centralSystem.SetSecurityHandler(handler)
	centralSystem.SetSecureFirmwareHandler(handler)
	centralSystem.SetLogHandler(handler)

	// Add handlers for dis/connection of charge points
	centralSystem.SetNewChargePointHandler(func(chargePoint ocpp16.ChargePointConnection) {
		handler.chargePoints[chargePoint.ID()] = &ChargePointState{connectors: map[int]*ConnectorInfo{}, transactions: map[int]*TransactionInfo{}}
		log.WithField("client", chargePoint.ID()).Info("new charge point connected")
		go exampleRoutine(chargePoint.ID(), handler)
	})
	centralSystem.SetChargePointDisconnectedHandler(func(chargePoint ocpp16.ChargePointConnection) {
		log.WithField("client", chargePoint.ID()).Info("charge point disconnected")
		delete(handler.chargePoints, chargePoint.ID())
	})
	ocppj.SetLogger(log.WithField("logger", "ocppj"))
	ws.SetLogger(log.WithField("logger", "websocket"))
	// Run central system
	log.Infof("starting central system on port %v", listenPort)
	centralSystem.Start(listenPort, "/{ws}")
	log.Info("stopped central system")
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	// Set this to DebugLevel if you want to retrieve verbose logs from the ocppj and websocket layers
	log.SetLevel(logrus.ErrorLevel)
}

# OCPP 1.6 Security Extension

Security extension for OCPP 1.6 is a set of additional features that can be added to the OCPP 1.6 protocol.
Additional handlers on both the central system and charge point side must be implemented in order to support these
features.

## Central System

To add support for security extension in the central system, you have the following handlers:

```go
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

```

## Charge Point

To add support for security extension in the charge point, you have the following handlers:

```go
handler := &ChargePointHandler{}
// Support callbacks for all OCPP 1.6 profiles
chargePoint.SetCoreHandler(handler)
chargePoint.SetFirmwareManagementHandler(handler)
chargePoint.SetLocalAuthListHandler(handler)
chargePoint.SetReservationHandler(handler)
chargePoint.SetRemoteTriggerHandler(handler)
chargePoint.SetSmartChargingHandler(handler)
// OCPP 1.6j Security extension
chargePoint.SetCertificateHandler(handler)
chargePoint.SetLogHandler(handler)
chargePoint.SetSecureFirmwareHandler(handler)
chargePoint.SetExtendedTriggerMessageHandler(handler)
chargePoint.SetSecurityHandler(handler)

```

## Additional remarks

### HTTP Basic Auth

The security extension specifies how to secure the communication between charge points and central systems
using HTTP Basic Auth and/or certificates. These are already provided in the websocket server/client
implementation.

Example charge point:

```go
wsClient := ws.NewClient()
wsClient.SetBasicAuth("foo", "bar")
cp := ocpp16.NewChargePoint(chargePointID, nil, wsClient)
```

Example central system:

```go
server := ws.NewServer()
server.SetBasicAuthHandler(func (chargePointID string, username string, password string) bool {
// todo Handle basic auth
return true
})
cs := ocpp16.NewCentralSystem(nil, server)
```

### Certificate-based authentication (mTLS)

The security extension specifies how to secure the communication between charge points and central systems
using mTLS (client certificates). The library provides the necessary functionality to configure TLS,
but mTLS itself is not in scope and should be handled by the user.

### Additional configuration keys

The OCPP 1.6 security extension introduces additional configuration keys.
These are not a part of the standard library, but they impact how the charge point should behave.

The charge point websocket client should be restarted when the `AuthorizationKey` configuration changes.


// Package data contains the data transfer functional block for OCPP 2.1.
// It enables parties to add custom commands and extensions to OCPP 2.1.
package data

import "github.com/lorenzodonini/ocpp-go/ocpp"

// ProfileName is the name of the data functional block.
const ProfileName = "data"

// ChargingStationHandler contains callbacks for handling incoming CSMS requests.
type ChargingStationHandler interface {
	// OnDataTransfer handles a DataTransferRequest from the CSMS.
	OnDataTransfer(request *DataTransferRequest) (response *DataTransferResponse, err error)
}

// CSMSHandler contains callbacks for handling incoming Charging Station requests.
type CSMSHandler interface {
	// OnDataTransfer handles a DataTransferRequest from a Charging Station.
	OnDataTransfer(chargingStationID string, request *DataTransferRequest) (response *DataTransferResponse, err error)
}

// Profile returns the data profile with all features.
func Profile() *ocpp.Profile {
	return ocpp.NewProfile(
		ProfileName,
		DataTransferFeature{},
	)
}

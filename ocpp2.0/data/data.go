// The data transfer functional block enables parties to add custom commands and extensions to OCPP 2.0.
package data

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Data transfer profile.
type CSMSHandler interface {
	// OnDataTransfer is called on the CSMS whenever a DataTransferRequest is received from a charging station.
	OnDataTransfer(chargingStationID string, request *DataTransferRequest) (confirmation *DataTransferResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Data transfer profile.
type ChargingStationHandler interface {
	// OnDataTransfer is called on a charging station whenever a DataTransferRequest is received from the CSMS.
	OnDataTransfer(request *DataTransferRequest) (confirmation *DataTransferResponse, err error)
}

const ProfileName = "data"

var Profile = ocpp.NewProfile(
	ProfileName,
	DataTransferFeature{},
	)

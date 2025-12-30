package main

import (
	"encoding/json"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/data"
)

type DataSample struct {
	SampleString string  `json:"sample_string"`
	SampleValue  float64 `json:"sample_value"`
}

func (handler *ChargingStationHandler) OnDataTransfer(request *data.DataTransferRequest) (response *data.DataTransferResponse, err error) {
	var dataSample DataSample
	err = json.Unmarshal(request.Data.([]byte), &dataSample)
	if err != nil {
		logDefault(request.GetFeatureName()).
			Errorf("invalid data received: %v", request.Data)
		return nil, err
	}
	logDefault(request.GetFeatureName()).
		Infof("data received: %v, %v", dataSample.SampleString, dataSample.SampleValue)
	return data.NewDataTransferResponse(data.DataTransferStatusAccepted), nil
}

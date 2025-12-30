package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/firmware"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/types"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
	"io"
	"net/http"
	"os"
	"time"
)

func (handler *ChargingStationHandler) OnPublishFirmware(request *firmware.PublishFirmwareRequest) (response *firmware.PublishFirmwareResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnUnpublishFirmware(request *firmware.UnpublishFirmwareRequest) (response *firmware.UnpublishFirmwareResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnUpdateFirmware(request *firmware.UpdateFirmwareRequest) (response *firmware.UpdateFirmwareResponse, err error) {
	retries := 0
	retryInterval := 30
	if request.Retries != nil {
		retries = *request.Retries
	}
	if request.RetryInterval != nil {
		retryInterval = *request.RetryInterval
	}
	logDefault(request.GetFeatureName()).Infof("starting update firmware procedure")
	go updateFirmware(request.Firmware.Location, request.Firmware.RetrieveDateTime, request.Firmware.InstallDateTime, retries, retryInterval)
	return firmware.NewUpdateFirmwareResponse(firmware.UpdateFirmwareStatusAccepted), nil
}

func updateFirmwareStatus(status firmware.FirmwareStatus, props ...func(request *firmware.FirmwareStatusNotificationRequest)) {
	statusConfirmation, err := chargingStation.FirmwareStatusNotification(status, props...)
	checkError(err)
	logDefault(statusConfirmation.GetFeatureName()).Infof("firmware status updated to %v", status)
}

// Retrieve data and install date are ignored for this test function.
func updateFirmware(location string, retrieveDate *types.DateTime, installDate *types.DateTime, retries int, retryInterval int) {
	updateFirmwareStatus(firmware.FirmwareStatusDownloading)
	err := downloadFile("/tmp/out.bin", location)
	if err != nil {
		logDefault(firmware.UpdateFirmwareFeatureName).Errorf("error while downloading file %v", err)
		updateFirmwareStatus(firmware.FirmwareStatusDownloadFailed)
		return
	}
	updateFirmwareStatus(firmware.FirmwareStatusDownloaded)
	// Simulate installation
	updateFirmwareStatus(firmware.FirmwareStatusInstalling)
	time.Sleep(time.Second * 5)
	// Notify completion
	updateFirmwareStatus(firmware.FirmwareStatusInstalled)
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/display"
)

func (c *CSMSHandler) OnNotifyDisplayMessages(chargingStationID string, request *display.NotifyDisplayMessagesRequest) (response *display.NotifyDisplayMessagesResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("received display messages for request %v:\n", request.RequestID)
	for _, msg := range request.MessageInfo {
		logDefault(chargingStationID, request.GetFeatureName()).Printf("%v", msg)
	}
	response = display.NewNotifyDisplayMessagesResponse()
	return
}

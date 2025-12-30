package main

import "github.com/xBlaz3kx/ocpp-go/ocpp2.1/display"

func (handler *ChargingStationHandler) OnClearDisplay(request *display.ClearDisplayRequest) (response *display.ClearDisplayResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("cleared display message %v", request.ID)
	response = display.NewClearDisplayResponse(display.ClearMessageStatusAccepted)
	return
}

func (handler *ChargingStationHandler) OnGetDisplayMessages(request *display.GetDisplayMessagesRequest) (response *display.GetDisplayMessagesResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("request %v to send display messages ignored", request.RequestID)
	response = display.NewGetDisplayMessagesResponse(display.MessageStatusUnknown)
	return
}

func (handler *ChargingStationHandler) OnSetDisplayMessage(request *display.SetDisplayMessageRequest) (response *display.SetDisplayMessageResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("accepted request to display message %v: %v", request.Message.ID, request.Message.Message.Content)
	response = display.NewSetDisplayMessageResponse(display.DisplayMessageStatusAccepted)
	return
}

package core

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"reflect"
)

// -------------------- Authorize --------------------
type AuthorizeRequest struct {
	IdTag string				`json:"idTag" validate:"required,max=20"`
}

type AuthorizeConfirmation struct {
	IdTagInfo ocpp.IdTagInfo	`json:"idTagInfo" validate:"required"`
}

type AuthorizeFeature struct {}

func (f AuthorizeFeature) GetFeatureName() string {
	return v16.AuthorizeFeatureName
}

func (f AuthorizeFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(AuthorizeRequest{})
}

func (f AuthorizeFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(AuthorizeConfirmation{})
}

func (r AuthorizeRequest) GetFeatureName() string {
	return v16.AuthorizeFeatureName
}

func (c AuthorizeConfirmation) GetFeatureName() string {
	return v16.AuthorizeFeatureName
}

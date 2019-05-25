package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
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
	return AuthorizeFeatureName
}

func (f AuthorizeFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(AuthorizeRequest{})
}

func (f AuthorizeFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(AuthorizeConfirmation{})
}

func (r AuthorizeRequest) GetFeatureName() string {
	return AuthorizeFeatureName
}

func (c AuthorizeConfirmation) GetFeatureName() string {
	return AuthorizeFeatureName
}

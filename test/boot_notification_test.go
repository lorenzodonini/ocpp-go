package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CoreTestSuite struct {
	suite.Suite
}

func (suite *CoreTestSuite) SetupTest() {
	coreProfile := ocpp.Profile{Features: make(map[string]ocpp.Feature)}
	feature := v16.BootNotificationFeature{}
	coreProfile.AddFeature(feature)
	ocpp.AddProfile(&coreProfile)
}

func GetBootNotificationRequest(t* testing.T, request ocpp.Request) *v16.BootNotificationRequest {
	assert.NotNil(t, request)
	result := request.(*v16.BootNotificationRequest)
	assert.NotNil(t, result)
	assert.IsType(t, v16.BootNotificationRequest{}, result)
	return result
}

func (suite *CoreTestSuite) TestBootNotificationValid() {
	t := suite.T()
	dataJson := `[2,"1234","BootNotification",{"chargePointModel": "model1", "chargePointVendor": "ABL"}]`
	call := ParseCall(dataJson, t)
	CheckCall(call, t, v16.BootNotificationFeatureName, "1234")
	request := GetBootNotificationRequest(t, call.Payload)
	assert.Equal(t, "model1", request.ChargePointModel)
	assert.Equal(t, "ABL", request.ChargePointVendor)
}

func (suite *CoreTestSuite) TestBootNotificationInvalid() {

}

func (suite *CoreTestSuite) TestBootNotificationMessage() {

}

func TestBootNotification(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}

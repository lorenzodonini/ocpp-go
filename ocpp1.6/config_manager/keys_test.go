package ocpp_16_config_manager

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/stretchr/testify/suite"
)

type keyTestSuite struct {
	suite.Suite
}

func (s *keyTestSuite) TestGetMandatoryKeysForProfile_Core() {
	keys := GetMandatoryKeysForProfile(core.ProfileName)

	s.Assert().ElementsMatch(keys, MandatoryCoreKeys)
}

func (s *keyTestSuite) TestGetMandatoryKeysForProfile_LocalAuth() {
	keys := GetMandatoryKeysForProfile(localauth.ProfileName)

	s.Assert().ElementsMatch(keys, MandatoryLocalAuthKeys)
}

func (s *keyTestSuite) TestGetMandatoryKeysForProfile_Mix() {
	keys := GetMandatoryKeysForProfile(core.ProfileName, localauth.ProfileName, firmware.ProfileName, smartcharging.ProfileName)

	expectedKeys := append(MandatoryCoreKeys, MandatoryLocalAuthKeys...)
	expectedKeys = append(expectedKeys, MandatoryFirmwareKeys...)
	expectedKeys = append(expectedKeys, MandatorySmartChargingKeys...)

	s.Assert().ElementsMatch(keys, expectedKeys)
}

func (s *keyTestSuite) TestGetMandatoryKeysForProfile_None() {
	keys := GetMandatoryKeysForProfile()
	s.Assert().Empty(keys)
}

func TestGetMandatoryKeysForProfile(t *testing.T) {
	suite.Run(t, new(keyTestSuite))
}

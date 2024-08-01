package ocpp_16_config_manager

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
)

var (
	val1 = "60"
	val2 = "ABCD"
)

type OcppConfigTest struct {
	suite.Suite
	config Config
}

func (s *OcppConfigTest) SetupTest() {
	s.config = Config{
		Version: 1,
		Keys: []core.ConfigurationKey{
			{
				Key:      "HeartbeatInterval",
				Readonly: false,
				Value:    &val1,
			}, {
				Key:      "ChargingScheduleAllowedChargingRateUnit",
				Readonly: true,
				Value:    &val2,
			}, {
				Key:      "AuthorizationCacheEnabled",
				Readonly: false,
				Value:    nil,
			},
		},
	}
}

func (s *OcppConfigTest) TestGetConfig() {
	s.Assert().Equal([]core.ConfigurationKey{
		{
			Key:      "HeartbeatInterval",
			Readonly: false,
			Value:    &val1,
		}, {
			Key:      "ChargingScheduleAllowedChargingRateUnit",
			Readonly: true,
			Value:    &val2,
		}, {
			Key:      "AuthorizationCacheEnabled",
			Readonly: false,
			Value:    nil,
		},
	}, s.config.GetConfig())

	// Overwrite the config
	s.config = Config{
		Version: 1,
		Keys:    []core.ConfigurationKey{},
	}

	s.Assert().Equal([]core.ConfigurationKey{}, s.config.GetConfig())
}

func (s *OcppConfigTest) TestGetConfigurationValue() {
	// Ok case
	value, err := s.config.GetConfigurationValue(HeartbeatInterval.String())
	s.Require().NoError(err)
	s.Assert().EqualValues("60", *value)

	// Invalid key
	value, err = s.config.GetConfigurationValue("Test4")
	s.Assert().Error(err)
	s.Assert().Nil(value)
}

func (s *OcppConfigTest) TestGetVersion() {
	s.Assert().EqualValues(1, s.config.GetVersion())

	s.config.Version = 1234

	s.Assert().EqualValues(1234, s.config.GetVersion())
}

func (s *OcppConfigTest) TestSetVersion() {
	s.config.SetVersion(1234)
	s.Assert().EqualValues(1234, s.config.Version)

	s.config.SetVersion(1)
	s.Assert().EqualValues(1, s.config.Version)
}

func (s *OcppConfigTest) TestUpdateKey() {
	// Ok case
	newVal := "1234"
	err := s.config.UpdateKey("HeartbeatInterval", &newVal)
	s.Assert().NoError(err)
	value, err := s.config.GetConfigurationValue("HeartbeatInterval")
	s.Require().NoError(err)
	s.Assert().EqualValues("1234", *value)

	// Invalid key
	err = s.config.UpdateKey("Test4", nil)
	s.Assert().Error(err)

	// Key cannot be updated
	err = s.config.UpdateKey("ChargingScheduleAllowedChargingRateUnit", nil)
	s.Assert().Error(err)

	// Check if the value was not updated
	value, err = s.config.GetConfigurationValue("ChargingScheduleAllowedChargingRateUnit")
	s.Assert().NoError(err)
	s.Assert().EqualValues("ABCD", *value)
}

func (s *OcppConfigTest) TestUpdateKeyReadability() {
	// Ok case
	err := s.config.UpdateKeyReadability(HeartbeatInterval.String(), true)
	s.Assert().NoError(err)

	configKey, isFound := lo.Find(s.config.Keys, func(item core.ConfigurationKey) bool {
		return item.Key == HeartbeatInterval.String()
	})
	s.Assert().True(isFound)
	s.Assert().EqualValues(true, configKey.Readonly)

	// Invalid key
	err = s.config.UpdateKeyReadability("Test4", true)
	s.Assert().Error(err)
}

func (s *OcppConfigTest) TestValidate() {
	s.config = NewEmptyConfiguration()
	s.config.Keys = DefaultCoreConfiguration()

	// Ok case
	err := s.config.Validate(MandatoryCoreKeys)
	s.Assert().NoError(err)

	// Missing mandatory key
	s.config.Keys = s.config.Keys[:len(s.config.Keys)-2]
	err = s.config.Validate(MandatoryCoreKeys)
	s.Assert().Error(err)
}

func TestOCPPConfig(t *testing.T) {
	suite.Run(t, new(OcppConfigTest))
}

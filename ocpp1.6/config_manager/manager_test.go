package ocpp_16_config_manager

import (
	"strings"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
)

type ConfigurationManagerTestSuite struct {
	suite.Suite
	manager *ManagerV16
}

func (s *ConfigurationManagerTestSuite) SetupTest() {
	profiles := []string{core.ProfileName, smartcharging.ProfileName, localauth.ProfileName}
	configuration, err := DefaultConfigurationFromProfiles(profiles...)
	s.Require().NoError(err)

	s.manager, err = NewV16ConfigurationManager(*configuration, profiles...)
	s.Assert().NoError(err)
}

func (s *ConfigurationManagerTestSuite) TestNewV16ConfigurationManager() {
	configuration, err := DefaultConfigurationFromProfiles(core.ProfileName)
	s.Assert().NoError(err)

	// Valid configuration
	manager, err := NewV16ConfigurationManager(*configuration, core.ProfileName)
	s.Assert().NoError(err)
	s.Assert().NotNil(manager)

	// Invalid configuration
	manager, err = NewV16ConfigurationManager(NewEmptyConfiguration(), core.ProfileName)
	s.Assert().Error(err)
	s.Assert().Nil(manager)
}

func (s *ConfigurationManagerTestSuite) TestGetConfiguration() {
	// todo
}

func (s *ConfigurationManagerTestSuite) TestUpdateConfiguration() {
	// Key found
	err := s.manager.UpdateKey(HeartbeatInterval, lo.ToPtr("123"))
	s.Assert().NoError(err)

	value, err := s.manager.GetConfigurationValue(HeartbeatInterval)
	s.Assert().NoError(err)
	s.Assert().NotNil(value)
	s.Assert().Equal("123", *value)

	err = s.manager.UpdateKey(HeartbeatInterval, nil)
	s.Assert().NoError(err)

	value, err = s.manager.GetConfigurationValue(HeartbeatInterval)
	s.Assert().NoError(err)
	s.Assert().Nil(value)

	// Key not found
	err = s.manager.UpdateKey("ExampleKey", lo.ToPtr("exampleValue"))
	s.Assert().Error(err)

	err = s.manager.UpdateKey("", lo.ToPtr("exampleValue"))
	s.Assert().Error(err)

	err = s.manager.UpdateKey("", nil)
	s.Assert().Error(err)
}

func (s *ConfigurationManagerTestSuite) TestGetConfigurationValue() {
	// Tested in the UpdateConfiguration test
}

func (s *ConfigurationManagerTestSuite) TestOnUpdateKey() {
	numExecutions := 0

	err := s.manager.OnUpdateKey(HeartbeatInterval, func(value *string) error {
		numExecutions++
		return nil
	})
	_ = s.manager.UpdateKey(HeartbeatInterval, lo.ToPtr("exampleValue"))
	_ = s.manager.UpdateKey(HeartbeatInterval, lo.ToPtr("exampleValue"))
	_ = s.manager.UpdateKey(HeartbeatInterval, lo.ToPtr("exampleValue"))
	s.Assert().NoError(err)
	s.Assert().Equal(3, numExecutions)

	err = s.manager.OnUpdateKey("ExampleKey", func(value *string) error {
		return nil
	})
	s.Assert().Error(err)

	err = s.manager.UpdateKey("ExampleKey", lo.ToPtr("exampleValue"))
	s.Assert().Error(err)

	err = s.manager.OnUpdateKey("", func(value *string) error {
		numExecutions++
		return nil
	})
	s.Assert().Error(err)

	err = s.manager.UpdateKey("", lo.ToPtr("exampleValue"))
	s.Assert().Error(err)
	s.Assert().Equal(3, numExecutions)
}

func (s *ConfigurationManagerTestSuite) TestSetConfiguration() {
	// todo
}

func (s *ConfigurationManagerTestSuite) TestValidateKey() {

	s.manager.RegisterCustomKeyValidator(func(key Key, value *string) bool {
		switch key {
		case HeartbeatInterval:
			if value == nil {
				return false
			}
			return true
		case LocalAuthListEnabled:
			if value == nil {
				return false
			}

			if strings.ToUpper(*value) == "TRUE" || strings.ToUpper(*value) == "FALSE" {
				return true
			}

			return false
		default:
			return false
		}
	})

	err := s.manager.ValidateKey(HeartbeatInterval, lo.ToPtr("123"))
	s.Assert().NoError(err)
	// Should fail - invalid value
	err = s.manager.ValidateKey(HeartbeatInterval, nil)
	s.Assert().Error(err)

	err = s.manager.ValidateKey(LocalAuthListEnabled, lo.ToPtr("true"))
	s.Assert().NoError(err)

	err = s.manager.ValidateKey(LocalAuthListEnabled, lo.ToPtr("false"))
	s.Assert().NoError(err)

	// Should fail - invalid value
	err = s.manager.ValidateKey(LocalAuthListEnabled, nil)
	s.Assert().Error(err)

	err = s.manager.ValidateKey(LocalAuthListEnabled, lo.ToPtr("aaaaa"))
	s.Assert().Error(err)

	// Should fail - invalid key
	err = s.manager.ValidateKey("ABCD", lo.ToPtr("aaaaa"))
	s.Assert().Error(err)

	err = s.manager.ValidateKey("ABCD", nil)
	s.Assert().Error(err)
}

func (s *ConfigurationManagerTestSuite) TestRegisterCustomKeyValidator() {
	exampleKey := Key("ExampleKey")
	numExecutions := 0

	s.manager.RegisterCustomKeyValidator(func(key Key, value *string) bool {
		numExecutions++
		return exampleKey == key && value != nil && *value == "exampleValue"
	})

	_ = s.manager.UpdateKey(exampleKey, lo.ToPtr("exampleValue"))
	s.Assert().Equal(1, numExecutions)
}

func (s *ConfigurationManagerTestSuite) TestGetMandatoryKeys() {
	configuration, err := DefaultConfigurationFromProfiles(core.ProfileName)
	s.Assert().NoError(err)

	// Valid configuration
	manager, err := NewV16ConfigurationManager(*configuration, core.ProfileName)
	s.Assert().NoError(err)

	s.ElementsMatch(MandatoryCoreKeys, manager.GetMandatoryKeys())

	configuration, err = DefaultConfigurationFromProfiles(core.ProfileName, localauth.ProfileName)
	s.Assert().NoError(err)
	manager, err = NewV16ConfigurationManager(*configuration, core.ProfileName, localauth.ProfileName)
	s.Assert().NoError(err)

	keys := append(MandatoryCoreKeys, MandatoryLocalAuthKeys...)
	s.ElementsMatch(keys, manager.GetMandatoryKeys())
}

func (s *ConfigurationManagerTestSuite) TestSetMandatoryKeys() {
	configuration, err := DefaultConfigurationFromProfiles(core.ProfileName)
	s.Assert().NoError(err)

	// Valid configuration
	manager, err := NewV16ConfigurationManager(*configuration, core.ProfileName)
	s.Assert().NoError(err)
	s.Assert().NotNil(manager)

}

func TestConfigurationManager(t *testing.T) {
	suite.Run(t, new(ConfigurationManagerTestSuite))
}

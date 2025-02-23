package ocpp_16_config_manager

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/stretchr/testify/suite"
)

type defaultsTestSuite struct {
	suite.Suite
}

func (suite *defaultsTestSuite) TestNewEmptyConfiguration() {
	config := NewEmptyConfiguration()
	suite.Equal(1, config.Version)
	suite.Empty(config.Keys)
}

func (suite *defaultsTestSuite) TestDefaultConfiguration() {

	// Default configuration with core, localauth and smartcharging profiles
	keys := append(DefaultCoreConfiguration(), DefaultLocalAuthConfiguration()...)
	keys = append(keys, DefaultSmartChargingConfiguration()...)

	config, err := DefaultConfigurationFromProfiles(core.ProfileName, localauth.ProfileName, smartcharging.ProfileName)
	suite.NoError(err)
	suite.Equal(1, config.Version)
	suite.NotEmpty(config.Keys)
	suite.ElementsMatch(keys, config.Keys)

	// Default configuration with unknown profile
	config, err = DefaultConfigurationFromProfiles(core.ProfileName, localauth.ProfileName, smartcharging.ProfileName, "unknown")
	suite.Error(err)
	suite.Nil(config)

	config, err = DefaultConfigurationFromProfiles()
	suite.Error(err)
	suite.Nil(config)
}

func (suite *defaultsTestSuite) TestDefaultCoreConfiguration() {
	config := DefaultCoreConfiguration()
	suite.NotEmpty(config)
}

func (suite *defaultsTestSuite) TestDefaultLocalAuthConfiguration() {
	config := DefaultLocalAuthConfiguration()
	suite.NotEmpty(config)
}

func (suite *defaultsTestSuite) TestDefaultSmartChargingConfiguration() {
	config := DefaultSmartChargingConfiguration()
	suite.NotEmpty(config)
}

func (suite *defaultsTestSuite) TestDefaultFirmwareConfiguration() {
	config := DefaultFirmwareConfiguration()
	suite.NotEmpty(config)
}

func TestDefaultConfigurations(t *testing.T) {
	suite.Run(t, new(defaultsTestSuite))
}

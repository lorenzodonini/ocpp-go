package ocpp_16_config_manager

import (
	"fmt"
	"strings"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/samber/lo"
)

func NewEmptyConfiguration() Config {
	return Config{
		Version: 1,
		Keys:    []core.ConfigurationKey{},
	}
}

func DefaultConfigurationFromProfiles(profiles ...string) (*Config, error) {
	keys := []core.ConfigurationKey{}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("no profiles provided")
	}

	for _, profile := range profiles {
		switch profile {
		case core.ProfileName:
			keys = append(keys, DefaultCoreConfiguration()...)
		case localauth.ProfileName:
			keys = append(keys, DefaultLocalAuthConfiguration()...)
		case smartcharging.ProfileName:
			keys = append(keys, DefaultSmartChargingConfiguration()...)
		case firmware.ProfileName:
			keys = append(keys, DefaultFirmwareConfiguration()...)
		default:
			return nil, fmt.Errorf("unknown profile %v", profile)
		}
	}

	return &Config{
		Version: 1,
		Keys:    keys,
	}, nil
}

func DefaultCoreConfiguration() []core.ConfigurationKey {
	return []core.ConfigurationKey{
		{
			Key:      AuthorizeRemoteTxRequests.String(),
			Readonly: false,
			Value:    lo.ToPtr("true"),
		},
		{
			Key:      ClockAlignedDataInterval.String(),
			Readonly: false,
			Value:    lo.ToPtr("0"),
		},
		{
			Key:      ConnectionTimeOut.String(),
			Readonly: false,
			Value:    lo.ToPtr("60"),
		},
		{
			Key:      GetConfigurationMaxKeys.String(),
			Readonly: false,
			Value:    lo.ToPtr("100"),
		},
		{
			Key:      HeartbeatInterval.String(),
			Readonly: false,
			Value:    lo.ToPtr("60"),
		},
		{
			Key:      LocalPreAuthorize.String(),
			Readonly: false,
			Value:    lo.ToPtr("false"),
		},
		{
			Key:      MeterValuesAlignedData.String(),
			Readonly: false,
			Value:    lo.ToPtr("true"),
		},
		{
			Key:      MeterValuesSampledData.String(),
			Readonly: false,
			Value: lo.ToPtr(strings.Join([]string{
				string(types.MeasurandVoltage),
				string(types.MeasurandCurrentImport),
				string(types.MeasurandPowerActiveImport),
				string(types.MeasurandEnergyActiveImportInterval),
				string(types.MeasueandSoC),
			}, ",")),
		},
		{
			Key:      MeterValueSampleInterval.String(),
			Readonly: false,
			Value:    lo.ToPtr("20"),
		},
		{
			Key:      NumberOfConnectors.String(),
			Readonly: true,
			Value:    lo.ToPtr("1"),
		},
		{
			Key:      ResetRetries.String(),
			Readonly: false,
			Value:    lo.ToPtr("3"),
		},
		{
			Key:      ConnectorPhaseRotation.String(),
			Readonly: true,
			Value:    lo.ToPtr("Unknown"),
		},
		{
			Key:      StopTransactionOnEVSideDisconnect.String(),
			Readonly: false,
			Value:    lo.ToPtr("true"),
		},
		{
			Key:      StopTransactionOnInvalidId.String(),
			Readonly: false,
			Value:    lo.ToPtr("true"),
		},
		{
			Key:      StopTxnAlignedData.String(),
			Readonly: false,
			Value: lo.ToPtr(strings.Join([]string{
				string(types.MeasurandVoltage),
				string(types.MeasurandCurrentImport),
				string(types.MeasurandPowerActiveImport),
				string(types.MeasurandEnergyActiveImportInterval),
				string(types.MeasueandSoC),
			}, ",")),
		},
		{
			Key:      StopTxnSampledData.String(),
			Readonly: false,
			Value: lo.ToPtr(strings.Join([]string{
				string(types.MeasurandVoltage),
				string(types.MeasurandCurrentImport),
				string(types.MeasurandPowerActiveImport),
				string(types.MeasurandEnergyActiveImportInterval),
				string(types.MeasueandSoC),
			}, ",")),
		},
		{
			Key:      SupportedFeatureProfiles.String(),
			Readonly: true,
			Value:    lo.ToPtr("Core"),
		},
		{
			Key:      TransactionMessageAttempts.String(),
			Readonly: false,
			Value:    lo.ToPtr("3"),
		},
		{
			Key:      TransactionMessageRetryInterval.String(),
			Readonly: false,
			Value:    lo.ToPtr("30"),
		},
		{
			Key:      UnlockConnectorOnEVSideDisconnect.String(),
			Readonly: false,
			Value:    lo.ToPtr("true"),
		},
	}
}

func DefaultLocalAuthConfiguration() []core.ConfigurationKey {
	return []core.ConfigurationKey{
		{
			Key:      LocalAuthListEnabled.String(),
			Readonly: false,
			Value:    lo.ToPtr("true"),
		},
		{
			Key:      LocalAuthListMaxLength.String(),
			Readonly: true,
			Value:    lo.ToPtr("100"),
		},
		{
			Key:      SendLocalListMaxLength.String(),
			Readonly: true,
			Value:    lo.ToPtr("100"),
		},
	}
}

func DefaultSmartChargingConfiguration() []core.ConfigurationKey {
	return []core.ConfigurationKey{
		{
			Key:      ChargeProfileMaxStackLevel.String(),
			Readonly: true,
			Value:    lo.ToPtr("5"),
		},
		{
			Key:      ChargingScheduleAllowedChargingRateUnit.String(),
			Readonly: true,
			Value:    lo.ToPtr("Current,Power"),
		},
		{
			Key:      ChargingScheduleMaxPeriods.String(),
			Readonly: true,
			Value:    lo.ToPtr("6"),
		},
		{
			Key:      MaxChargingProfilesInstalled.String(),
			Readonly: true,
			Value:    lo.ToPtr("5"),
		},
	}
}

func DefaultFirmwareConfiguration() []core.ConfigurationKey {
	return []core.ConfigurationKey{
		{
			Key:      SupportedFileTransferProtocols.String(),
			Readonly: true,
			Value:    lo.ToPtr("HTTP,HTTPS,FTP,FTPS,SFTP"),
		},
	}
}

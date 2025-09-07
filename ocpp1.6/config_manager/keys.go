package ocpp_16_config_manager

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
)

const (
	/* ----------------- Core keys ----------------------- */

	AllowOfflineTxForUnknownId        = Key("AllowOfflineTxForUnknownId")
	AuthorizationCacheEnabled         = Key("AuthorizationCacheEnabled")
	AuthorizeRemoteTxRequests         = Key("AuthorizeRemoteTxRequests")
	BlinkRepeat                       = Key("BlinkRepeat")
	ClockAlignedDataInterval          = Key("ClockAlignedDataInterval")
	ConnectionTimeOut                 = Key("ConnectionTimeOut")
	GetConfigurationMaxKeys           = Key("GetConfigurationMaxKeys")
	HeartbeatInterval                 = Key("HeartbeatInterval")
	LightIntensity                    = Key("LightIntensity")
	LocalAuthorizeOffline             = Key("LocalAuthorizeOffline")
	LocalPreAuthorize                 = Key("LocalPreAuthorize")
	MaxEnergyOnInvalidId              = Key("MaxEnergyOnInvalidId")
	MeterValuesAlignedData            = Key("MeterValuesAlignedData")
	MeterValuesAlignedDataMaxLength   = Key("MeterValuesAlignedDataMaxLength")
	MeterValuesSampledData            = Key("MeterValuesSampledData")
	MeterValuesSampledDataMaxLength   = Key("MeterValuesSampledDataMaxLength")
	MeterValueSampleInterval          = Key("MeterValueSampleInterval")
	MinimumStatusDuration             = Key("MinimumStatusDuration")
	NumberOfConnectors                = Key("NumberOfConnectors")
	ResetRetries                      = Key("ResetRetries")
	ConnectorPhaseRotation            = Key("ConnectorPhaseRotation")
	ConnectorPhaseRotationMaxLength   = Key("ConnectorPhaseRotationMaxLength")
	StopTransactionOnEVSideDisconnect = Key("StopTransactionOnEVSideDisconnect")
	StopTransactionOnInvalidId        = Key("StopTransactionOnInvalidId")
	StopTxnAlignedData                = Key("StopTxnAlignedData")
	StopTxnAlignedDataMaxLength       = Key("StopTxnAlignedDataMaxLength")
	StopTxnSampledData                = Key("StopTxnSampledData")
	StopTxnSampledDataMaxLength       = Key("StopTxnSampledDataMaxLength")
	SupportedFeatureProfiles          = Key("SupportedFeatureProfiles")
	SupportedFeatureProfilesMaxLength = Key("SupportedFeatureProfilesMaxLength")
	TransactionMessageAttempts        = Key("TransactionMessageAttempts")
	TransactionMessageRetryInterval   = Key("TransactionMessageRetryInterval")
	UnlockConnectorOnEVSideDisconnect = Key("UnlockConnectorOnEVSideDisconnect")
	WebSocketPingInterval             = Key("WebSocketPingInterval")

	/* ----------------- LocalAuthList keys ----------------------- */

	LocalAuthListEnabled   = Key("LocalAuthListEnabled")
	LocalAuthListMaxLength = Key("LocalAuthListMaxLength")
	SendLocalListMaxLength = Key("SendLocalListMaxLength")

	/* ----------------- Reservation keys ----------------------- */

	ReserveConnectorZeroSupported = Key("ReserveConnectorZeroSupported")

	/* ----------------- Firmware keys ----------------------- */

	SupportedFileTransferProtocols = Key("SupportedFileTransferProtocols")

	/* ----------------- SmartCharging keys ----------------------- */

	ChargeProfileMaxStackLevel              = Key("ChargeProfileMaxStackLevel")
	ChargingScheduleAllowedChargingRateUnit = Key("ChargingScheduleAllowedChargingRateUnit")
	ChargingScheduleMaxPeriods              = Key("ChargingScheduleMaxPeriods")
	MaxChargingProfilesInstalled            = Key("MaxChargingProfilesInstalled")
	ConnectorSwitch3to1PhaseSupported       = Key("ConnectorSwitch3to1PhaseSupported")

	/* ----------------- ISO15118 keys ----------------------- */
	CentralContractValidationAllowed = Key("CentralContractValidationAllowed")
	CertificateSignedMaxChainSize    = Key("CertificateSignedMaxChainSize")
	CertSigningWaitMinimum           = Key("CertSigningWaitMinimum")
	CertSigningRepeatTimes           = Key("CertSigningRepeatTimes")
	CertificateStoreMaxLength        = Key("CertificateStoreMaxLength")
	ContractValidationOffline        = Key("ContractValidationOffline")
	ISO15118PnCEnabled               = Key("ISO15118PnCEnabled")

	/* ----------------- Security extension keys ----------------------- */
	AuthorizationData              = Key("AuthorizationData")
	AdditionalRootCertificateCheck = Key("AdditionalRootCertificateCheck")
	CpoName                        = Key("CpoName")
	SecurityProfile                = Key("SecurityProfile")
)

var (
	MandatoryCoreKeys = []Key{
		AuthorizeRemoteTxRequests,
		ClockAlignedDataInterval,
		ConnectionTimeOut,
		GetConfigurationMaxKeys,
		HeartbeatInterval,
		LocalPreAuthorize,
		MeterValuesAlignedData,
		MeterValuesSampledData,
		MeterValueSampleInterval,
		NumberOfConnectors,
		ResetRetries,
		ConnectorPhaseRotation,
		StopTransactionOnEVSideDisconnect,
		StopTransactionOnInvalidId,
		StopTxnAlignedData,
		StopTxnSampledData,
		SupportedFeatureProfiles,
		TransactionMessageAttempts,
		TransactionMessageRetryInterval,
		UnlockConnectorOnEVSideDisconnect,
	}

	MandatoryLocalAuthKeys = []Key{
		LocalAuthListEnabled,
		LocalAuthListMaxLength,
		SendLocalListMaxLength,
	}

	MandatorySmartChargingKeys = []Key{
		MaxChargingProfilesInstalled,
		ChargingScheduleMaxPeriods,
		ChargingScheduleAllowedChargingRateUnit,
		ChargeProfileMaxStackLevel,
	}

	MandatoryFirmwareKeys = []Key{
		SupportedFileTransferProtocols,
	}

	MandatoryISO15118Keys = []Key{
		ISO15118PnCEnabled,
		ContractValidationOffline,
	}

	// Security extension does not have any mandatory keys
)

func GetMandatoryKeysForProfile(profiles ...string) []Key {
	mandatoryKeys := []Key{}

	for _, profile := range profiles {
		switch profile {
		case core.ProfileName:
			mandatoryKeys = append(mandatoryKeys, MandatoryCoreKeys...)
		case smartcharging.ProfileName:
			mandatoryKeys = append(mandatoryKeys, MandatorySmartChargingKeys...)
		case localauth.ProfileName:
			mandatoryKeys = append(mandatoryKeys, MandatoryLocalAuthKeys...)
		case firmware.ProfileName:
			mandatoryKeys = append(mandatoryKeys, MandatoryFirmwareKeys...)
			// todo IS15118 mandatory keys validation
		}
	}

	return mandatoryKeys
}

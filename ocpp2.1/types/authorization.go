package types

import "gopkg.in/go-playground/validator.v9"

type AuthorizationStatus string

const (
	AuthorizationStatusAccepted           AuthorizationStatus = "Accepted"
	AuthorizationStatusBlocked            AuthorizationStatus = "Blocked"
	AuthorizationStatusExpired            AuthorizationStatus = "Expired"
	AuthorizationStatusInvalid            AuthorizationStatus = "Invalid"
	AuthorizationStatusConcurrentTx       AuthorizationStatus = "ConcurrentTx"
	AuthorizationStatusNoCredit           AuthorizationStatus = "NoCredit"
	AuthorizationStatusNotAllowedTypeEVSE AuthorizationStatus = "NotAllowedTypeEVSE"
	AuthorizationStatusNotAtThisLocation  AuthorizationStatus = "NotAtThisLocation"
	AuthorizationStatusNotAtThisTime      AuthorizationStatus = "NotAtThisTime"
	AuthorizationStatusUnknown            AuthorizationStatus = "Unknown"
)

func isValidAuthorizationStatus(fl validator.FieldLevel) bool {
	status := AuthorizationStatus(fl.Field().String())
	switch status {
	case AuthorizationStatusAccepted, AuthorizationStatusBlocked, AuthorizationStatusExpired, AuthorizationStatusInvalid, AuthorizationStatusConcurrentTx, AuthorizationStatusNoCredit, AuthorizationStatusNotAllowedTypeEVSE, AuthorizationStatusNotAtThisLocation, AuthorizationStatusNotAtThisTime, AuthorizationStatusUnknown:
		return true
	default:
		return false
	}
}

// ID Token
type IdTokenType string

const (
	IdTokenTypeCentral         IdTokenType = "Central"
	IdTokenTypeEMAID           IdTokenType = "eMAID"
	IdTokenTypeISO14443        IdTokenType = "ISO14443"
	IdTokenTypeISO15693        IdTokenType = "ISO15693"
	IdTokenTypeKeyCode         IdTokenType = "KeyCode"
	IdTokenTypeLocal           IdTokenType = "Local"
	IdTokenTypeMacAddress      IdTokenType = "MacAddress"
	IdTokenTypeNoAuthorization IdTokenType = "NoAuthorization"
)

func isValidIdTokenType(fl validator.FieldLevel) bool {
	tokenType := IdTokenType(fl.Field().String())
	switch tokenType {
	case IdTokenTypeCentral, IdTokenTypeEMAID, IdTokenTypeISO14443, IdTokenTypeISO15693, IdTokenTypeKeyCode, IdTokenTypeLocal, IdTokenTypeMacAddress, IdTokenTypeNoAuthorization:
		return true
	default:
		return false
	}
}

func isValidIdToken(sl validator.StructLevel) {
	idToken := sl.Current().Interface().(IdToken)
	// validate required idToken value except `NoAuthorization` type
	switch idToken.Type {
	case IdTokenTypeCentral, IdTokenTypeEMAID, IdTokenTypeISO14443, IdTokenTypeISO15693, IdTokenTypeKeyCode, IdTokenTypeLocal, IdTokenTypeMacAddress:
		if idToken.IdToken == "" {
			sl.ReportError(idToken.IdToken, "IdToken", "IdToken", "required", "")
		}
	}
}

type AdditionalInfo struct {
	AdditionalIdToken string `json:"additionalIdToken" validate:"required,max=36"`
	Type              string `json:"type" validate:"required,max=50"`
}

type IdToken struct {
	IdToken        string           `json:"idToken" validate:"max=255"`
	Type           IdTokenType      `json:"type" validate:"required,idTokenType,max=20"`
	AdditionalInfo []AdditionalInfo `json:"additionalInfo,omitempty" validate:"omitempty,dive"`
}

type GroupIdToken struct {
	IdToken string      `json:"idToken" validate:"max=36"`
	Type    IdTokenType `json:"type" validate:"required,idTokenType"`
}

func isValidGroupIdToken(sl validator.StructLevel) {
	groupIdToken := sl.Current().Interface().(GroupIdToken)
	// validate required idToken value except `NoAuthorization` type
	switch groupIdToken.Type {
	case IdTokenTypeCentral, IdTokenTypeEMAID, IdTokenTypeISO14443, IdTokenTypeISO15693, IdTokenTypeKeyCode, IdTokenTypeLocal, IdTokenTypeMacAddress:
		if groupIdToken.IdToken == "" {
			sl.ReportError(groupIdToken.IdToken, "IdToken", "IdToken", "required", "")
		}
	}
}

type IdTokenInfo struct {
	Status              AuthorizationStatus `json:"status" validate:"required,authorizationStatus21"`
	CacheExpiryDateTime *DateTime           `json:"cacheExpiryDateTime,omitempty" validate:"omitempty"`
	ChargingPriority    int                 `json:"chargingPriority,omitempty" validate:"min=-9,max=9"`
	Language1           string              `json:"language1,omitempty" validate:"max=8"`
	Language2           string              `json:"language2,omitempty" validate:"max=8"`
	GroupIdToken        *GroupIdToken       `json:"groupIdToken,omitempty"`
	PersonalMessage     *MessageContent     `json:"personalMessage,omitempty"`
}

// NewIdTokenInfo creates an IdTokenInfo. Optional parameters may be set afterwards on the initialized struct.
func NewIdTokenInfo(status AuthorizationStatus) *IdTokenInfo {
	return &IdTokenInfo{Status: status}
}

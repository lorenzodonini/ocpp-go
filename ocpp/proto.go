package ocpp

type Validatable interface {
	validate() error
}

type Feature interface {
	FeatureName() string
}

type Request struct {
	Validatable
	Feature Feature
}

type Confirmation struct {
	Validatable
	Feature Feature
}

type MessageType int

const (
	CALL 		= 2
	CALL_RESULT = 3
	CALL_ERROR 	= 4
)

type Message struct {
	MessageTypeId MessageType	`json:"messageTypeId"`
	UniqueId string 			`json:"uniqueId"`	//Max 36 chars
	Validatable
}

func (m* Message) validate() error {
	return nil
}

type Call struct {
	Message
	Action string				`json:"action"`
	Payload Request				`json:"payload"`
}

type CallResult struct {
	Message
	Payload Confirmation		`json:"payload"`
}

type CallError struct {
	Message
	ErrorCode ErrorCode			`json:"errorCode"`
	ErrorDescription string		`json:"errorDescription"`
	ErrorDetails interface{}	`json:"errorDetails"`
}
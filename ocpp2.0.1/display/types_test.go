package display

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type displayTestSuite struct {
	suite.Suite
}

func (suite *displayTestSuite) TestMessageInfoValidation() {
	var testTable = []tests.GenericTestEntry{
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), EndDateTime: types.NewDateTime(time.Now().Add(1 * time.Hour)), TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}, Display: &types.Component{Name: "name1"}}, true},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), EndDateTime: types.NewDateTime(time.Now().Add(1 * time.Hour)), TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle}, false},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8}}, false},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: "invalidState", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{MessageInfo{ID: 42, State: MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{MessageInfo{ID: 42, Priority: "invalidPriority", State: MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{MessageInfo{ID: -1, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, TransactionID: ">36..................................", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), EndDateTime: types.NewDateTime(time.Now().Add(1 * time.Hour)), TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}, Display: &types.Component{}}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), testTable)
}

func TestDisplaySuite(t *testing.T) {
	suite.Run(t, new(displayTestSuite))
}

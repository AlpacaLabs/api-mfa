package entities

import (
	eventV1 "github.com/AlpacaLabs/protorepo-event-go/alpacalabs/event/v1"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"github.com/golang/protobuf/proto"
	"github.com/rs/xid"
)

type SendEvent struct {
	eventV1.EventInfo
	eventV1.TraceInfo
	Sent     bool
	Catalyst mfaV1.DeliverCodeRequest
	Payload  proto.Message
}

func NewSendEvent(traceInfo eventV1.TraceInfo, catalyst mfaV1.DeliverCodeRequest, payload proto.Message) SendEvent {
	return SendEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo: traceInfo,
		Sent:      false,
		Catalyst:  catalyst,
		Payload:   payload,
	}
}

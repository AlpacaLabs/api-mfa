package services

import (
	"context"

	"github.com/sfreiberg/gotwilio"
)

type SendSmsInput struct {
	TwilioClient      *gotwilio.Twilio
	TwilioPhoneNumber string
	To                string
	Message           string
}

func (s *Service) SendSms(ctx context.Context, in SendSmsInput) error {
	_, exception, err := in.TwilioClient.SendSMS(in.TwilioPhoneNumber, in.To, in.Message, "", "")
	if err != nil {
		return err
	}
	if exception != nil {
		return exception
	}
	return nil
}

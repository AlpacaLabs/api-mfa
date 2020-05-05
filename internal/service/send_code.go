package service

import (
	"context"
	"errors"
	"fmt"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	"github.com/AlpacaLabs/mfa/internal/db"
	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

// SendCode lets the client choose which email address or phone number to deliver the code to.
func (s *Service) SendCode(ctx context.Context, request *mfaV1.SendCodeRequest) (*mfaV1.SendCodeResponse, error) {

	// TODO verify requester ID == request.AccountId

	// primary key for the code entity
	codeID := request.CodeId
	accountID := request.AccountId

	// Dial the Account Service
	accountConn, err := grpc.Dial(s.config.AccountGRPCAddress)
	if err != nil {
		return nil, err
	}
	accountClient := accountV1.NewAccountServiceClient(accountConn)

	// Get the account's info
	res, err := accountClient.GetAccount(ctx, &accountV1.GetAccountRequest{
		AccountIdentifier: &accountV1.GetAccountRequest_AccountId{AccountId: accountID},
	})
	if err != nil {
		return nil, err
	}

	// Get the value of the email or phone number the client chose
	var emailAddress, phoneNumber string
	if eid := request.GetEmailAddressId(); eid != "" {
		for _, e := range res.Account.EmailAddresses {
			if e.Confirmed && e.Id == eid {
				emailAddress = e.EmailAddress
				break
			}
		}
	} else if pid := request.GetPhoneNumberId(); pid != "" {
		for _, p := range res.Account.PhoneNumbers {
			if p.Confirmed && p.Id == eid {
				phoneNumber = p.PhoneNumber
				break
			}
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "must provide either email_address_id or phone_number_id to deliver MFA code")
	}

	// This shouldn't ever happen, but it's included just in case...
	if emailAddress == "" && phoneNumber == "" {
		return nil, errors.New("must provide email_address_id or phone_number_id that reference something extant")
	}

	var out *mfaV1.SendCodeResponse

	err = s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// Retrieve the MFA code
		code, err := tx.GetCode(ctx, codeID)
		if err != nil {
			return err
		}

		// TODO use transactional outbox pattern instead

		hermesConn, err := grpc.Dial(s.config.HermesGRPCAddress)
		if err != nil {
			return err
		}

		if emailAddress != "" {
			// Send MFA code via email
			emailClient := hermesV1.NewSendEmailServiceClient(hermesConn)
			_, err = emailClient.SendEmail(ctx, &hermesV1.SendEmailRequest{
				Email: &hermesV1.Email{
					Subject: "Multi-factor Authentication Code",
					Body: &hermesV1.Body{
						Name: "",
						Intros: []string{
							fmt.Sprintf("Your MFA code is: %s", code.Code),
						},
						Actions:   nil,
						Outros:    nil,
						Greeting:  "",
						Signature: "",
					},
					To: []*hermesV1.Recipient{
						{
							Email: emailAddress,
						},
					},
				},
			})
		} else {
			// Send MFA code via SMS
			smsClient := hermesV1.NewSendSmsServiceClient(hermesConn)
			_, err = smsClient.SendSms(ctx, &hermesV1.SendSmsRequest{
				To:      phoneNumber,
				Message: fmt.Sprintf("Your MFA code is: %s", code.Code),
			})
		}

		if err != nil {
			return err
		}

		out = &mfaV1.SendCodeResponse{}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

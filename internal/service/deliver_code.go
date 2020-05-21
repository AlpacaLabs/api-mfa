package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlpacaLabs/api-mfa/internal/db/entities"
	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"
	"github.com/golang/protobuf/proto"

	log "github.com/sirupsen/logrus"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"

	"github.com/AlpacaLabs/api-mfa/internal/db"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

var (
	ErrUnsupportedCodeDeliveryMethod = errors.New("must provide id for email or phone number to receive code")
)

// DeliverCode lets the client choose which email address or phone number to deliver the code to.
func (s Service) DeliverCode(ctx context.Context, request *mfaV1.DeliverCodeRequest) (*mfaV1.DeliverCodeResponse, error) {
	funcName := "DeliverCode"

	// primary key for the code entity
	codeID := request.CodeId

	accountClient := accountV1.NewAccountServiceClient(s.accountConn)

	var account *accountV1.Account
	var emailAddressID string
	var phoneNumberID string
	if pid := request.GetPhoneNumberId(); pid != "" {
		phoneNumberID = pid
		// Get the account's info by phone number ID
		res, err := accountClient.GetAccount(ctx, &accountV1.GetAccountRequest{
			AccountIdentifier: &accountV1.GetAccountRequest_PhoneNumberId{PhoneNumberId: pid},
		})
		if err != nil {
			return nil, err
		}

		// Set the account
		account = res.Account
	} else if eid := request.GetEmailAddressId(); eid != "" {
		emailAddressID = eid
		// Get the account's info by email address ID
		res, err := accountClient.GetAccount(ctx, &accountV1.GetAccountRequest{
			AccountIdentifier: &accountV1.GetAccountRequest_EmailAddressId{EmailAddressId: eid},
		})
		if err != nil {
			return nil, err
		}

		// Set the account
		account = res.Account
	} else {
		return nil, ErrUnsupportedCodeDeliveryMethod
	}

	log.Infof("looked up account %s", account.Id)

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// Retrieve the MFA code
		mfaCode, err := tx.GetCode(ctx, codeID)
		if err != nil {
			return err
		}

		log.Infof("Looked up MFA code: %s %s", mfaCode.Id, mfaCode.Code)

		// TODO verify requester ID == mfaCode.AccountId

		var payload proto.Message

		var transactionalOutboxTable string

		if emailAddressID != "" {
			payload = s.buildSendEmailRequest()
			transactionalOutboxTable = db.TableForSendEmailRequest
		} else if phoneNumberID != "" {
			payload = s.buildSendSmsRequest()
			transactionalOutboxTable = db.TableForSendSmsRequest
		}

		// Create the event entity that will be persisted to the transactional outbox
		event, err := entities.NewEvent(ctx, request, payload)
		if err != nil {
			return fmt.Errorf("failed to create event in %s: %w", funcName, err)
		}

		// Persist the event to the transactional outbox
		return tx.CreateEvent(ctx, event, transactionalOutboxTable)
	})

	if err != nil {
		return nil, err
	}

	return &mfaV1.DeliverCodeResponse{}, nil
}

func (s Service) buildSendEmailRequest() *hermesV1.SendEmailRequest {
	// TODO build email
	return &hermesV1.SendEmailRequest{
		Email: nil,
	}
}

func (s Service) buildSendSmsRequest() *hermesV1.SendSmsRequest {
	// TODO build sms message
	return &hermesV1.SendSmsRequest{
		To:      "",
		Message: "",
	}
}

package service

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"

	"github.com/AlpacaLabs/api-mfa/internal/db"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

var (
	ErrUnsupportedCodeDeliveryMethod = errors.New("must provide id for email or phone number to receive code")
)

// DeliverCode lets the client choose which email address or phone number to deliver the code to.
func (s *Service) DeliverCode(ctx context.Context, request *mfaV1.DeliverCodeRequest) (*mfaV1.DeliverCodeResponse, error) {

	// primary key for the code entity
	codeID := request.CodeId

	accountClient := accountV1.NewAccountServiceClient(s.accountConn)

	var account *accountV1.Account
	if pid := request.GetPhoneNumberId(); pid != "" {
		// Look up phone number
		phoneNumber, err := accountClient.GetPhoneNumberByID(ctx, &accountV1.GetPhoneNumberByIDRequest{
			Id: pid,
		})
		if err != nil {
			return nil, err
		}

		// Get the account's info
		res, err := accountClient.GetAccount(ctx, &accountV1.GetAccountRequest{
			// TODO support look up by phone number ID so we only have to make 1 grpc call
			AccountIdentifier: &accountV1.GetAccountRequest_PhoneNumber{PhoneNumber: phoneNumber.PhoneNumber.PhoneNumber},
		})
		if err != nil {
			return nil, err
		}

		// Set the account
		account = res.Account
	} else if eid := request.GetEmailAddressId(); eid != "" {
		// Look up email address
		emailAddress, err := accountClient.GetEmailAddressByID(ctx, &accountV1.GetEmailAddressByIDRequest{
			Id: eid,
		})
		if err != nil {
			return nil, err
		}

		// Get the account's info
		res, err := accountClient.GetAccount(ctx, &accountV1.GetAccountRequest{
			// TODO support look up by email address ID so we only have to make 1 grpc call
			AccountIdentifier: &accountV1.GetAccountRequest_EmailAddress{EmailAddress: emailAddress.EmailAddress.EmailAddress},
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

	var out *mfaV1.DeliverCodeResponse

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// Retrieve the MFA code
		mfaCode, err := tx.GetCode(ctx, codeID)
		if err != nil {
			return err
		}

		log.Infof("Looked up MFA code: %s %s", mfaCode.Id, mfaCode.Code)

		// TODO verify requester ID == mfaCode.AccountId

		// Write to transactional outbox so the code can get delivered later
		if err := tx.CreateTxobForCode(ctx, *request); err != nil {
			return err
		}

		out = &mfaV1.DeliverCodeResponse{}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

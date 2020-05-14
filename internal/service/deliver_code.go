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

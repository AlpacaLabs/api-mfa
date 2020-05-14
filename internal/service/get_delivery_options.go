package service

import (
	"context"

	"github.com/AlpacaLabs/api-mfa/internal/db"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

// GetDeliveryOptions will initiate the MFA flow by providing an account ID.
// It will persist a new MFA code to the database and return to the client
// a set of possible authentication options, including any confirmed email addresses
// and phone numbers belonging to the account.
func (s *Service) GetDeliveryOptions(ctx context.Context, request *mfaV1.GetDeliveryOptionsRequest) (*mfaV1.GetDeliveryOptionsResponse, error) {
	// The account that is trying to authenticate w/ MFA
	accountID := request.AccountId

	// Create a new code
	code := newCode(accountID)

	// Call the Accounts service to get their email addresses and phone numbers
	codeOptions, err := s.getCodeOptions(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Persist the code to the DB
	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// Persist the code to the DB
		if err := tx.CreateCode(ctx, code); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &mfaV1.GetDeliveryOptionsResponse{
		CodeId:      code.Id,
		CodeOptions: codeOptions,
	}, nil
}

func (s *Service) getCodeOptions(ctx context.Context, accountID string) (*mfaV1.CodeDeliveryOptions, error) {
	var emailAddresses []*accountV1.EmailAddress
	var phoneNumbers []*accountV1.PhoneNumber

	accountClient := accountV1.NewAccountServiceClient(s.accountConn)
	res, err := accountClient.GetAccount(ctx, &accountV1.GetAccountRequest{
		AccountIdentifier: &accountV1.GetAccountRequest_AccountId{AccountId: accountID},
	})
	if err != nil {
		return nil, err
	}

	for _, e := range res.Account.EmailAddresses {
		if e.Confirmed {
			emailAddresses = append(emailAddresses, e)
		}
	}

	for _, p := range res.Account.PhoneNumbers {
		if p.Confirmed {
			phoneNumbers = append(phoneNumbers, p)
		}
	}

	codeOptions := &mfaV1.CodeDeliveryOptions{
		EmailAddresses: []*mfaV1.EmailAddressOption{},
		PhoneNumbers:   []*mfaV1.PhoneNumberOption{},
	}

	// TODO consider masking phone numbers and emails

	for _, e := range emailAddresses {
		codeOptions.EmailAddresses = append(codeOptions.EmailAddresses, &mfaV1.EmailAddressOption{
			Id:           e.Id,
			AccountId:    e.AccountId,
			EmailAddress: e.EmailAddress,
		})
	}

	for _, p := range phoneNumbers {
		codeOptions.PhoneNumbers = append(codeOptions.PhoneNumbers, &mfaV1.PhoneNumberOption{
			Id:          p.Id,
			AccountId:   p.AccountId,
			PhoneNumber: p.PhoneNumber,
		})
	}

	return codeOptions, nil
}

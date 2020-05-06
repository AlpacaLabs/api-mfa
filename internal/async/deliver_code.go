package async

import (
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"google.golang.org/grpc"
)

type deliverCodeInput struct {
	hermesConn     *grpc.ClientConn
	code           *mfaV1.MFACode
	emailAddressID string
	phoneNumberID  string
}

//func deliverCode(in deliverCodeInput) {
//	ctx := context.TODO()
//
//	hermesConn := in.hermesConn
//	mfaCode := in.code
//	emailAddressID := in.emailAddressID
//	phoneNumberID := in.phoneNumberID
//
//	// TODO look up email or phone number from Account service, depending on which is provided
//
//	if emailAddressID != "" {
//		// Send MFA code via email
//		emailClient := hermesV1.NewSendEmailServiceClient(hermesConn)
//		_, err := emailClient.SendEmail(ctx, &hermesV1.SendEmailRequest{
//			Email: &hermesV1.Email{
//				Subject: "Multi-factor Authentication Code",
//				Body: &hermesV1.Body{
//					Name: "",
//					Intros: []string{
//						fmt.Sprintf("Your MFA code is: %s", mfaCode.Code),
//					},
//					Actions:   nil,
//					Outros:    nil,
//					Greeting:  "",
//					Signature: "",
//				},
//				To: []*hermesV1.Recipient{
//					{
//						Email: emailAddress,
//					},
//				},
//			},
//		})
//		if err != nil {
//			// TODO how do we handle async errors?
//			return
//		}
//
//	} else {
//		// Send MFA code via SMS
//		smsClient := hermesV1.NewSendSmsServiceClient(hermesConn)
//		_, err := smsClient.SendSms(ctx, &hermesV1.SendSmsRequest{
//			To:      phoneNumber,
//			Message: fmt.Sprintf("Your MFA code is: %s", mfaCode.Code),
//		})
//
//		if err != nil {
//			// TODO how do we handle async errors?
//			return
//		}
//	}
//}

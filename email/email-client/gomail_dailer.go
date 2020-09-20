package email_client

// mockgen -source=email/email-client/gomail_dailer.go -destination=email/email-client/mocks/mocks_gomail_dailer.go -package=mocks
import "gopkg.in/gomail.v2"

type GomailDialer interface {
	DialAndSend(messages ...*gomail.Message) error
}

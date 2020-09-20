package email_client

import "gopkg.in/gomail.v2"

type GomailDialer interface {
	DialAndSend(messages ...*gomail.Message) error
}

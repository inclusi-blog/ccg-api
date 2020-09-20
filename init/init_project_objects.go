package init

import (
	"ccg-api/configuration"
	"ccg-api/controller"
	. "ccg-api/email/configuration"
	emailControllers "ccg-api/email/controller"
	email_client "ccg-api/email/email-client"
	"ccg-api/email/service"
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

var (
	healthController = controller.HealthController{}
	emailController emailControllers.EmailController
)

func Objects(configData *configuration.ConfigData) {
	emailClientConfig := NewEmailClientConfig(configData.Email)
	gomailDialer := buildGomailDialer(emailClientConfig)
	emailService := service.NewEmailService(email_client.NewEmailClient(emailClientConfig.TempDir(), gomailDialer), emailClientConfig)
	emailController = emailControllers.NewEmailController(emailService, emailClientConfig)
}

func buildGomailDialer(config EmailClientConfig) email_client.GomailDialer {
	gomailDialer := gomail.NewDialer(config.SmtpHost(), config.SmtpPort(), config.Username(), config.Password())
	gomailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify()}
	return gomailDialer
}

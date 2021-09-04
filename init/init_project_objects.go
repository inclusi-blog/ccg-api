package init

import (
	"ccg-api/configuration"
	"ccg-api/controller"
	. "ccg-api/email/configuration"
	emailControllers "ccg-api/email/controller"
	emailClient "ccg-api/email/email-client"
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
	emailService := service.NewEmailService(emailClient.NewEmailClient(emailClientConfig.TempDir(), gomailDialer), emailClientConfig)
	emailController = emailControllers.NewEmailController(emailService, emailClientConfig)
}

func buildGomailDialer(config EmailClientConfig) emailClient.GomailDialer {
	gomailDialer := gomail.NewDialer(config.SmtpHost(), config.SmtpPort(), config.Username(), config.Password())
	gomailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify()}
	return gomailDialer
}

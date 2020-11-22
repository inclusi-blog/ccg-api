package configuration

import (
	"ccg-api/configuration"
	"os"
)

type EmailClientConfig interface {
	SmtpHost() string
	SmtpPort() int
	Username() string
	Password() string
	InsecureSkipVerify() bool
	TempDir() string
	ValidGolaEmailDomain() []string
	DefaultGolaEmailSender() string
	UnsupportedAttachmentExtensions() []string
	PermissibleTotalSizeOfAttachments() int
	BaseTemplateFilePath() string
	LogoUrls() configuration.LogoUrls
	OtherUrls() configuration.Urls
}

type emailClientConfig struct {
	email configuration.Email
}

func NewEmailClientConfig(email configuration.Email) EmailClientConfig {
	return emailClientConfig{email}
}

func (config emailClientConfig) OtherUrls() configuration.Urls {
	return config.email.OtherUrls
}

func (config emailClientConfig) SmtpHost() string {
	return config.email.SmtpHost
}

func (config emailClientConfig) SmtpPort() int {
	return config.email.SmtpPort
}

func (config emailClientConfig) Username() string {
	return config.email.Username
}

func (config emailClientConfig) Password() string {
	return os.Getenv("SMTP_CLIENT_PASSWORD")
}

func (config emailClientConfig) InsecureSkipVerify() bool {
	return config.email.InsecureSkipVerify
}

func (config emailClientConfig) TempDir() string {
	return os.TempDir()
}

func (config emailClientConfig) ValidGolaEmailDomain() []string {
	return config.email.ValidMensuvadiEmailDomains
}

func (config emailClientConfig) DefaultGolaEmailSender() string {
	return config.email.DefaultMensuvadiEmailSender
}

func (config emailClientConfig) UnsupportedAttachmentExtensions() []string {
	return config.email.UnsupportedAttachmentExtensions
}

func (config emailClientConfig) PermissibleTotalSizeOfAttachments() int {
	return config.email.PermissibleAttachmentSizeInBytes
}

func (config emailClientConfig) BaseTemplateFilePath() string {
	return config.email.BaseTemplateFilePath
}

func (config emailClientConfig) LogoUrls() configuration.LogoUrls {
	return config.email.LogoUrls
}

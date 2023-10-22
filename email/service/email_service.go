package service

import (
	"bytes"
	"ccg-api/constants"
	"ccg-api/email/configuration"
	"ccg-api/email/email-client"
	"ccg-api/email/email-client/email_client_request"
	"ccg-api/email/models"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/inclusi-blog/gola-utils/mask_util"
	"html/template"
	"strings"
)

type EmailService interface {
	Send(ctx *gin.Context, email models.Email) *golaerror.Error
}

type emailService struct {
	emailClient email_client.EmailClient
	emailConfig configuration.EmailClientConfig
}

func NewEmailService(emailClient email_client.EmailClient, emailConfig configuration.EmailClientConfig) EmailService {
	return emailService{emailClient: emailClient, emailConfig: emailConfig}
}

func (emailService emailService) Send(ctx *gin.Context, email models.Email) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "EmailService").WithField("method", "Send")
	if email.IncludeBaseTemplate {
		var templateParseError error
		email.Body.Content, templateParseError = emailService.embedContentInBaseTemplate(ctx, email.Body.Content)
		if templateParseError != nil {
			logger.Error("Could not parse template ", templateParseError)
			return &constants.InternalServerError
		}
	}
	request := email_client_request.EmailClientRequest{
		From:        email.From,
		To:          email.To,
		Subject:     email.Subject,
		Body:        email.Body,
		Attachments: email.Attachments,
	}
	err := emailService.emailClient.Send(ctx, &request)
	if err != nil {
		logger.Error("Error received from email client ", err)
		return &constants.InternalServerError
	}

	var maskedEmail []string
	for _, mailId := range email.To {
		maskedEmail = append(maskedEmail, mask_util.MaskEmail(ctx, mailId))
	}
	logger.Infof("Email sent successfully to %s", strings.Join(maskedEmail, ", "))
	return nil
}

func (emailService emailService) embedContentInBaseTemplate(ctx *gin.Context, content string) (string, error) {
	logger := logging.GetLogger(ctx).WithField("class", "EmailService").WithField("method", "embedContentInBaseTemplate")
	contentBuffer := new(bytes.Buffer)
	var err error

	baseTemplate, err := template.New("content").Parse(content)
	if err != nil {
		logger.Error("Error while parsing content ", err)
		return "", err
	}

	// adding func to avoid escaping conditional HTML comments
	finalTemplate, err := baseTemplate.New("base").Funcs(template.FuncMap{
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}).ParseFiles(emailService.emailConfig.BaseTemplateFilePath())

	if err != nil {
		logger.Error("Error while parsing template files ", err)
		return "", err
	}

	fields := map[string]interface{}{
		"LogoUrl": emailService.emailConfig.LogoUrls(),
		"Urls":    emailService.emailConfig.OtherUrls(),
	}
	err = finalTemplate.ExecuteTemplate(contentBuffer, "base", fields)
	if err != nil {
		logger.Error("Error while executing template ", err)
		return "", err
	}
	return contentBuffer.String(), nil
}

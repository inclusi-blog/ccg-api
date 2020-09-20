package http_request_response

import (
	"ccg-api/email/models"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
)

type EmailRequest struct {
	From                string       `json:"from" binding:"required" validate:"email,validGolaEmailDomain" example:"abc@gola.xyz"`
	To                  []string     `json:"to" binding:"required" validate:"gt=0,dive,email" example:"abc@gmail.com"`
	Subject             string       `json:"subject" binding:"required" validate:"notblank" example:"base64 encoded value"`
	Body                MessageBody  `json:"message_body" binding:"required"`
	Attachments         []Attachment `json:"attachments" validate:"uniqueAttachments,totalAttachmentSizeWithinPermissibleLimit,dive"`
	IncludeBaseTemplate bool         `json:"include_base_template" example:"true"`
}

func (emailRequest EmailRequest) ToEmailModel(ctx *gin.Context) (models.Email, error) {
	messageBody, messageBodyError := emailRequest.buildMessageBody(emailRequest.Body)
	if messageBodyError != nil {
		logging.GetLogger(ctx).Error("Error while parsing message body")
		return models.Email{}, messageBodyError
	}

	attachments, attachmentError := emailRequest.buildAttachments(emailRequest.Attachments)
	if attachmentError != nil {
		logging.GetLogger(ctx).Error("Error while parsing email attachments")
		return models.Email{}, attachmentError
	}

	email := models.Email{
		From:                emailRequest.From,
		To:                  emailRequest.To,
		Subject:             emailRequest.Subject,
		Body:                messageBody,
		Attachments:         attachments,
		IncludeBaseTemplate: emailRequest.IncludeBaseTemplate,
	}
	return email, nil
}

func (emailRequest EmailRequest) buildAttachments(attachmentsInHTTPRequest []Attachment) ([]models.Attachment, error) {
	var attachments []models.Attachment
	for _, attachmentInHttpRequest := range attachmentsInHTTPRequest {
		attachment, attachmentError := emailRequest.buildAttachment(attachmentInHttpRequest)
		if attachmentError != nil {
			return []models.Attachment{}, attachmentError
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}

func (emailRequest EmailRequest) buildMessageBody(messageBodyInRequest MessageBody) (models.MessageBody, error) {
	decodedContent, err := messageBodyInRequest.GetDecodedContent()
	if err != nil {
		return models.MessageBody{}, err
	}

	messageBody := models.MessageBody{
		MimeType: messageBodyInRequest.MimeType,
		Content:  decodedContent,
	}

	return messageBody, nil
}

func (emailRequest EmailRequest) buildAttachment(attachmentInRequest Attachment) (models.Attachment, error) {
	decodedBytes, err := attachmentInRequest.GetDecodedData()
	if err != nil {
		return models.Attachment{}, err
	}
	attachment := models.Attachment{
		FileName: attachmentInRequest.FileName,
		Data:     decodedBytes,
	}
	return attachment, nil
}

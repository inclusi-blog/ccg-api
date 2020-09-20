package email_client_request

import (
	"ccg-api/email/models"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"path"
)

type EmailClientRequest struct {
	From        string
	To          []string
	Subject     string
	Body        models.MessageBody
	Attachments []models.Attachment
}

func (request EmailClientRequest) ToMessage(ctx *gin.Context, tempAttachmentDir string) (*gomail.Message, error) {
	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("From", request.From)
	gomailMessage.SetHeaders(map[string][]string{"To": request.To})

	gomailMessage.SetHeader("Subject", request.Subject)
	gomailMessage.SetBody(request.Body.MimeType, request.Body.Content)

	for _, attachment := range request.Attachments {
		if err := request.addAttachment(ctx, tempAttachmentDir, attachment, gomailMessage); err != nil {
			return nil, err
		}
	}
	return gomailMessage, nil
}

func (request EmailClientRequest) addAttachment(
	ctx *gin.Context,
	tempAttachmentDir string,
	attachment models.Attachment,
	gomailMessage *gomail.Message) error {

	filePath, fileCreationError := request.createFile(tempAttachmentDir, attachment)
	if fileCreationError != nil {
		logging.GetLogger(ctx).Error("Failed to create file to send attachment:", fileCreationError)
		return fileCreationError
	}
	gomailMessage.Attach(filePath)
	return nil
}

func (request EmailClientRequest) createFile(attachmentDir string, attachment models.Attachment) (string, error) {
	filePath := path.Join(attachmentDir, attachment.FileName)
	if err := ioutil.WriteFile(filePath, attachment.Data, 0644); err != nil {
		return "", err
	}
	return filePath, nil
}

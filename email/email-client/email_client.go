package email_client

import (
	"ccg-api/email/email-client/email_client_request"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"os"
	"path"
)

type EmailClient interface {
	Send(ctx *gin.Context, emailRequest *email_client_request.EmailClientRequest) (err error)
}

type emailClient struct {
	tempProcessingDir string
	gomailDialer      GomailDialer
}

func NewEmailClient(tempDir string, dialer GomailDialer) EmailClient {
	logger := logging.NewLoggerEntry()
	dirCreationError := os.MkdirAll(tempDir, 0755)
	if dirCreationError != nil {
		logger.Warn("Failed to create directory", dirCreationError)
	}
	client := emailClient{tempProcessingDir: tempDir, gomailDialer: dialer}
	return &client
}

func (emailClient *emailClient) Send(ctx *gin.Context, emailRequest *email_client_request.EmailClientRequest) (err error) {
	var tmpDirToUseForCurrentRequest = path.Join(emailClient.tempProcessingDir, uuid.New().String())
	if err := os.MkdirAll(tmpDirToUseForCurrentRequest, 0755); err != nil {
		logging.GetLogger(ctx).Error("Failed to create temp directory, required for e-mail attachments, directory: ", tmpDirToUseForCurrentRequest, err)
		return err
	}
	defer func() {
		err := os.RemoveAll(tmpDirToUseForCurrentRequest)
		if err != nil {
			logging.GetLogger(ctx).Warnf("Failed to delete temp attachment dir: %s, error: %s", tmpDirToUseForCurrentRequest, err)
		}
	}()

	message, _ := emailRequest.ToMessage(nil, tmpDirToUseForCurrentRequest)

	return emailClient.gomailDialer.DialAndSend(message)
}

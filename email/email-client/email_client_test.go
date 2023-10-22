package email_client

import (
	"ccg-api/email/email-client/email_client_request"
	mockemailclient "ccg-api/email/email-client/mocks"
	"ccg-api/email/email-client/test_helper"
	"ccg-api/email/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

type emailClientTestSuite struct {
	suite.Suite
	context *gin.Context

	mockCtrl     *gomock.Controller
	gomailDialer *mockemailclient.MockGomailDialer
}

func TestEmailClientTestSuite(t *testing.T) {
	suite.Run(t, new(emailClientTestSuite))
}

func (suite *emailClientTestSuite) SetupTest() {
	suite.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.gomailDialer = mockemailclient.NewMockGomailDialer(suite.mockCtrl)
}

func (suite *emailClientTestSuite) TestSend_ShouldReturnErrorIfUnableToSendEmailUsingGoClient() {
	request := email_client_request.EmailClientRequest{
		From:    "gola@gola.xyz",
		To:      []string{"someone@gmail.com"},
		Subject: "Hi...",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User...",
		},
	}

	tempDirForAttachingFiles := path.Join(os.TempDir(), "ccg-email-client-send-unit-test-dir-for-error-scenario")
	_ = os.MkdirAll(tempDirForAttachingFiles, 0755)
	defer os.RemoveAll(tempDirForAttachingFiles)

	emailClient := NewEmailClient(tempDirForAttachingFiles, suite.gomailDialer)
	logging.NewLoggerEntry().Debug("Temp dir for attachments: ", tempDirForAttachingFiles)

	suite.gomailDialer.EXPECT().DialAndSend(gomock.Any()).Return(errors.New("failed to connect SMTP server"))

	err := emailClient.Send(suite.context, &request)

	suite.NotNil(err)
}

func (suite *emailClientTestSuite) TestSend_ShouldSuccessfullySendEmailUsingGomailClient() {

	request := email_client_request.EmailClientRequest{
		From:    "gola@gola.xyz",
		To:      []string{"someone@gmail.com"},
		Subject: "Hi...",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User...",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment1.txt",
				Data:     []byte("Attachment1 Data"),
			},
			{
				FileName: "attachment2.txt",
				Data:     []byte("Attachment2 Data"),
			},
		},
	}

	tempDirForAttachingFiles := path.Join(os.TempDir(), "ccg-email-client-send-unit-test-dir-for-success-scenario")
	_ = os.MkdirAll(tempDirForAttachingFiles, 0755)
	defer os.RemoveAll(tempDirForAttachingFiles)

	emailClient := NewEmailClient(tempDirForAttachingFiles, suite.gomailDialer)
	logging.NewLoggerEntry().Debug("Temp dir for attachments: ", tempDirForAttachingFiles)

	var actualFromEmails []string
	var actualToEmails []string
	var actualSubjects []string
	var actualMessageCount int
	var sanitizedMessageContent string
	var contentsInTempDirWhileSendingEmail []os.FileInfo
	suite.gomailDialer.EXPECT().DialAndSend(gomock.Any()).Return(nil).Do(func(messages ...*gomail.Message) {
		actualMessageCount = len(messages)
		actualFromEmails = messages[0].GetHeader("From")
		actualToEmails = messages[0].GetHeader("To")
		actualSubjects = messages[0].GetHeader("Subject")
		sanitizedMessageContent = test_helper.SanitizeMessage(messages[0])
		contentsInTempDirWhileSendingEmail, _ = ioutil.ReadDir(tempDirForAttachingFiles)
	})

	err := emailClient.Send(suite.context, &request)

	suite.Nil(err)
	suite.Equal(actualMessageCount, 1)
	suite.Equal(actualFromEmails, []string{"gola@gola.xyz"})
	suite.Equal(actualToEmails, []string{"someone@gmail.com"})
	suite.Equal(actualSubjects, []string{"Hi..."})
	suite.Equal(sanitizedMessageContent, test_helper.SanitizeMessageContent(expectedMessageContentForEmailClientTest()))

	//Email client must create a random subdir within the temp directory before writing the contents to generate files for attachment
	for _, contentInTempDir := range contentsInTempDirWhileSendingEmail {
		suite.Equal(true, contentInTempDir.IsDir())
	}

	//Post completion of sending an email, the request specific temp directory should be cleaned by the email client
	contentsInTempDirPostSendingEmail, _ := ioutil.ReadDir(tempDirForAttachingFiles)
	suite.Less(len(contentsInTempDirPostSendingEmail), len(contentsInTempDirWhileSendingEmail))
}

func expectedMessageContentForEmailClientTest() string {
	return `Mime-Version: 1.0
Date: Sun, 23 Feb 2020 00:34:14 +0530
From: gola@gola.xyz
To: first@gmail.com, second@gmail.com
Subject: Hi...
Content-Type: multipart/mixed;
 boundary=5de10fdd5e00e6fadb59e85a1eb0e114db94cd6b1aa36c370687be1c1e57

--5de10fdd5e00e6fadb59e85a1eb0e114db94cd6b1aa36c370687be1c1e57
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain; charset=UTF-8

Hello User...
--5de10fdd5e00e6fadb59e85a1eb0e114db94cd6b1aa36c370687be1c1e57
Content-Disposition: attachment; filename="attachment1.txt"
Content-Transfer-Encoding: base64
Content-Type: application/octet-stream; name="attachment1.txt"

QXR0YWNobWVudDEgRGF0YQ==
--5de10fdd5e00e6fadb59e85a1eb0e114db94cd6b1aa36c370687be1c1e57
Content-Disposition: attachment; filename="attachment2.txt"
Content-Transfer-Encoding: base64
Content-Type: application/octet-stream; name="attachment2.txt"

QXR0YWNobWVudDIgRGF0YQ==
--5de10fdd5e00e6fadb59e85a1eb0e114db94cd6b1aa36c370687be1c1e57--
`
}

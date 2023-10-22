package email_client_request

import (
	"ccg-api/email/email-client/test_helper"
	"ccg-api/email/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

type emailClientRequestTestSuite struct {
	suite.Suite
	ctx *gin.Context
}

func TestMessageBodyTestSuite(t *testing.T) {
	suite.Run(t, new(emailClientRequestTestSuite))
}

func (suite *emailClientRequestTestSuite) SetupTest() {
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.ctx.Request, _ = http.NewRequest("GET", "", nil)
}

// noinspection GoNilness
func (suite *emailClientRequestTestSuite) TestToMessage_ShouldReturnExpectedGomailMessage() {
	emailClientRequest := EmailClientRequest{
		From:    "gola@gola.xyz",
		To:      []string{"first@gmail.com", "second@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
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

	tempAttachmentDir := path.Join(os.TempDir(), "email-send-unit-test-read-write-dir-for-gomail-message-creation")
	_ = os.Mkdir(tempAttachmentDir, 0777)
	defer os.RemoveAll(tempAttachmentDir)

	actualMessage, err := emailClientRequest.ToMessage(suite.ctx, tempAttachmentDir)
	suite.Nil(err)

	suite.Equal([]string{"gola@gola.xyz"}, actualMessage.GetHeader("From"))
	suite.Equal([]string{"first@gmail.com", "second@gmail.com"}, actualMessage.GetHeader("To"))
	suite.Equal([]string{"Hi!"}, actualMessage.GetHeader("Subject"))

	expectedMsg := test_helper.SanitizeMessageContent(expectedMessageContentForEmailClientRequestTest())
	actualMsg := test_helper.SanitizeMessage(actualMessage)
	suite.Equal(expectedMsg, actualMsg)
}

func (suite *emailClientRequestTestSuite) TestToMessage_ShouldReturnErrorIfUnableToBuildMessageDueToNonExistenceOfTempAttachmentDir() {
	emailClientRequest := EmailClientRequest{
		From:    "gola@gola.xyz",
		To:      []string{"first@gmail.com", "second@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment.txt",
				Data:     []byte("Attachment Data"),
			},
		},
	}

	tempAttachmentDir := path.Join(os.TempDir(), "email-send-unit-test-non-existent-dir")
	defer os.RemoveAll(tempAttachmentDir)

	_, err := emailClientRequest.ToMessage(suite.ctx, tempAttachmentDir)

	suite.NotNil(err)
}

func expectedMessageContentForEmailClientRequestTest() string {
	return `Mime-Version: 1.0
Date: Sun, 23 Feb 2020 00:34:14 +0530
From: gola@gola.xyz
To: first@gmail.com, second@gmail.com
Subject: Hi!
Content-Type: multipart/mixed;
 boundary=5de10fdd5e00e6fadb59e85a1eb0e114db94cd6b1aa36c370687be1c1e57

--5de10fdd5e00e6fadb59e85a1eb0e114db94cd6b1aa36c370687be1c1e57
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain; charset=UTF-8

Hello User!
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

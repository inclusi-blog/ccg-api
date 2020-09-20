package http_request_response

import (
	"ccg-api/email/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type emailRequestTestSuite struct {
	suite.Suite
	context  *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestEmailRequestTestSuite(t *testing.T) {
	suite.Run(t, new(emailRequestTestSuite))
}

func (suite *emailRequestTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "some-url", nil)
}

func (suite *emailRequestTestSuite) TestToEmail_ShouldReturnValidEmailModelObject() {
	emailRequest := EmailRequest{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: MessageBody{
			Content: "TWVzc2FnZSBCb2R5IQ==",
		},
		Attachments: []Attachment{
			{
				FileName: "attachment1.txt",
				Data:     "QXR0YWNobWVudDEgRGF0YSE=",
			},
			{
				FileName: "attachment2.txt",
				Data:     "QXR0YWNobWVudDIgRGF0YSE=",
			},
		},
		IncludeBaseTemplate: true,
	}

	expectedEmailModel := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			Content: "Message Body!",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment1.txt",
				Data:     []byte("Attachment1 Data!"),
			},
			{
				FileName: "attachment2.txt",
				Data:     []byte("Attachment2 Data!"),
			},
		},
		IncludeBaseTemplate: true,
	}

	actualEmailModel, err := emailRequest.ToEmailModel(suite.context)
	suite.Equal(expectedEmailModel, actualEmailModel)
	suite.Nil(err)
}

func (suite *emailRequestTestSuite) TestGetDecodedData_ShouldReturnErrorIfUnableToDecodeMessageBody() {
	emailRequest := EmailRequest{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: MessageBody{
			Content: "Plain text body",
		},
	}

	_, err := emailRequest.ToEmailModel(suite.context)
	suite.NotNil(err)
}

func (suite *emailRequestTestSuite) TestGetDecodedData_ShouldReturnErrorIfUnableToDecodeDataInAttachment() {
	emailRequest := EmailRequest{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: MessageBody{
			Content: "TWVzc2FnZSBCb2R5IQ==",
		},
		Attachments: []Attachment{
			{
				FileName: "attachment1.txt",
				Data:     "QXR0YWNobWVudDEgRGF0YSE=",
			},
			{
				FileName: "attachment2.txt",
				Data:     "Plain text data",
			},
		},
	}

	_, err := emailRequest.ToEmailModel(suite.context)
	suite.NotNil(err)
}

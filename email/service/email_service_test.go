package service

import (
	"ccg-api/configuration"
	"ccg-api/constants"
	"ccg-api/email/email-client/email_client_request"
	mockEmailClient "ccg-api/email/email-client/mocks"
	"ccg-api/email/mocks"
	"ccg-api/email/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type emailServiceTestSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	context      *gin.Context
	recorder     *httptest.ResponseRecorder
	emailClient  *mockEmailClient.MockEmailClient
	emailConfig  *mocks.MockEmailClientConfig
	emailService EmailService
}

func TestEmailServiceTestSuite(t *testing.T) {
	suite.Run(t, new(emailServiceTestSuite))
}

func (suite *emailServiceTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "some-url", nil)
	suite.emailClient = mockEmailClient.NewMockEmailClient(suite.mockCtrl)
	suite.emailConfig = mocks.NewMockEmailClientConfig(suite.mockCtrl)
	suite.emailService = NewEmailService(suite.emailClient, suite.emailConfig)
}

func (suite *emailServiceTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite emailServiceTestSuite) TestSendEmailShouldNotIncludeBaseTemplateAndSendEmailSuccessfully() {
	email := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment1.pdf",
				Data:     []byte("Attachement1 Data!"),
			},
		},
	}

	expectedEmailClientRequest := email_client_request.EmailClientRequest{
		From:        email.From,
		To:          email.To,
		Subject:     email.Subject,
		Body:        email.Body,
		Attachments: email.Attachments,
	}

	suite.emailClient.EXPECT().Send(suite.context, &expectedEmailClientRequest)

	err := suite.emailService.Send(suite.context, email)
	suite.Nil(err)
}

func (suite emailServiceTestSuite) TestSendEmailShouldIncludeBaseTemplateAndSendEmailSuccessfully() {
	email := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment1.pdf",
				Data:     []byte("Attachement1 Data!"),
			},
		},
		IncludeBaseTemplate: true,
	}

	suite.emailConfig.EXPECT().BaseTemplateFilePath().Return("../../email_templates/base_email_template.html")
	suite.emailConfig.EXPECT().LogoUrls().Return(configuration.LogoUrls{
		Mensuvadi:       "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		Facebook:        "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		Instagram:       "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		Twitter:         "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		LinkedIn:        "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		DownloadIOS:     "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		DownloadAndroid: "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
	})
	suite.emailConfig.EXPECT().OtherUrls().Return(configuration.Urls{
		HelpCenter:    "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		PrivacyPolicy: "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		Unsubscribe:   "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		FAQUrl:        "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
	})
	suite.emailClient.EXPECT().Send(suite.context, gomock.Any()).Do(func(ctx *gin.Context, request *email_client_request.EmailClientRequest) {
		suite.Equal(email.From, request.From)
		suite.Equal(email.To, request.To)
		suite.Equal(email.Subject, request.Subject)
		suite.Equal(email.Body.MimeType, request.Body.MimeType)
		suite.Equal(email.Attachments, request.Attachments)
	}).Return(nil).Times(1)

	err := suite.emailService.Send(suite.context, email)
	suite.Nil(err)
}

func (suite emailServiceTestSuite) TestSendEmailShouldSendEmailReturnErrorIfClientUnableToSendEmail() {
	email := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
	}

	errMsg := "failed to send email"
	suite.emailClient.EXPECT().Send(suite.context, gomock.Any()).Return(errors.New(errMsg))
	expectedError := &constants.InternalServerError

	actualError := suite.emailService.Send(suite.context, email)
	suite.Equal(expectedError, actualError)
}

func (suite emailServiceTestSuite) TestSendEmailShouldIncludeBaseTemplateAndThrowErrorIfTemplateNotFound() {
	email := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment1.pdf",
				Data:     []byte("Attachement1 Data!"),
			},
		},
		IncludeBaseTemplate: true,
	}

	suite.emailConfig.EXPECT().BaseTemplateFilePath().Return("base_email_template.html")

	err := suite.emailService.Send(suite.context, email)
	suite.Equal(&constants.InternalServerError, err)
}

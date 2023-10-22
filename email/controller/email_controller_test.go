package controller

import (
	"bytes"
	"ccg-api/constants"
	"ccg-api/email/http_request_response"
	"ccg-api/email/mocks"
	"ccg-api/email/models"
	"ccg-api/util"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type emailControllerTestSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	recorder     *httptest.ResponseRecorder
	context      *gin.Context
	emailService *mocks.MockEmailService
	controller   EmailController
	emailConfig  *mocks.MockEmailClientConfig
}

func TestEmailControllerTestSuite(t *testing.T) {
	suite.Run(t, new(emailControllerTestSuite))
}

const MaxPermissibleAttachmentSize = 9 * 1024 * 1024

func (suite *emailControllerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.emailService = mocks.NewMockEmailService(suite.mockCtrl)
	suite.emailConfig = mocks.NewMockEmailClientConfig(suite.mockCtrl)

	suite.emailConfig.EXPECT().SmtpHost().Return("smtp-host")
	suite.emailConfig.EXPECT().SmtpPort().Return(1234)
	suite.emailConfig.EXPECT().Username().Return("")
	suite.emailConfig.EXPECT().Password().Return("")
	suite.emailConfig.EXPECT().InsecureSkipVerify().Return(true)
	suite.emailConfig.EXPECT().TempDir().Return("/tmp")
	suite.emailConfig.EXPECT().ValidGolaEmailDomain().Return([]string{"gola.xyz"})
	suite.emailConfig.EXPECT().DefaultGolaEmailSender().Return("gola@gola.xyz")
	suite.emailConfig.EXPECT().PermissibleTotalSizeOfAttachments().Return(MaxPermissibleAttachmentSize)
	suite.emailConfig.EXPECT().UnsupportedAttachmentExtensions().Return([]string{"exe"})

	suite.controller = NewEmailController(suite.emailService, suite.emailConfig)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenFromEmailIsMissing() {
	request := suite.validEmailRequest()
	request.From = ""
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(constants.PayloadValidationErrorCode, response.ErrorCode)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenFromEmailIsNotValidEmailAddress() {
	request := suite.validEmailRequest()
	request.From = "invalid-email.com"
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(constants.PayloadValidationErrorCode, response.ErrorCode)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenFromEmailIsNotGolaDomain() {
	request := suite.validEmailRequest()
	request.From = "noreply@gmail.com"
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(constants.PayloadValidationErrorCode, response.ErrorCode)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenToEmailIsEmpty() {
	request := suite.validEmailRequest()
	request.To = []string{}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenToEmailIsInvalid() {
	request := suite.validEmailRequest()
	request.To = []string{"invalid-to-email.com"}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenOneOfToEmailIsInvalid() {
	request := suite.validEmailRequest()
	request.To = []string{"valid@gmail.com", "invalid-to-email.com"}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenSubjectIsEmpty() {
	request := suite.validEmailRequest()
	request.Subject = ""
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenSubjectIsBlank() {
	request := suite.validEmailRequest()
	request.Subject = "   "
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenMessageBodyIsEmpty() {
	request := suite.validEmailRequest()
	request.Body = http_request_response.MessageBody{}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenContentInMessageBodyIsMissing() {
	request := suite.validEmailRequest()
	request.Body = http_request_response.MessageBody{
		Content: "",
	}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenContentInMessageBodyIsBlank() {
	request := suite.validEmailRequest()
	request.Body = http_request_response.MessageBody{
		Content: base64.StdEncoding.EncodeToString([]byte("        ")),
	}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenMessageBodyIsNotInBase64EncodedFormat() {
	request := suite.validEmailRequest()
	request.Body = http_request_response.MessageBody{
		Content: "Message body without base64 encoding",
	}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenAttachmentsContainEmptyElement() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{{}}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenAttachmentsContainElementWithEmptyFileName() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{{
		Data: "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
	}}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenAttachmentDataIsNotBase64EncodedString() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{{
		FileName: "attachment.pdf",
		Data:     "Data is not in base64 encoded format",
	}}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenAttachmentFileExtensionIsNotSupported() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{{
		FileName: "attachment.exe",
		Data:     "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
	}}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenAttachmentFileStartsWithDot() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{{
		FileName: ".attachment",
		Data:     "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
	}}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenAttachmentFileHasNoExtension() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{{
		FileName: "attachment",
		Data:     "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
	}}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenAttachmentFileHasEmptyExtension() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{{
		FileName: "attachment.",
		Data:     "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
	}}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestWhenAttachmentHasFileNameButNoData() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{{
		FileName: "attachment.pdf",
		Data:     "",
	}}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestForDuplicateAttachments() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{
		{
			FileName: "attachment.pdf",
			Data:     "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
		},
		{
			FileName: "attachment.pdf",
			Data:     "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
		},
	}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(constants.PayloadValidationErrorCode, response.ErrorCode)
	suite.Equal("One or more of the request parameters are missing or invalid", response.ErrorMessage)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldThrowBadRequestIfAttachmentSizeIsGreaterThanAllowedLimit() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{
		{
			FileName: "attachment.pdf",
			Data:     base64.StdEncoding.EncodeToString([]byte(strings.Repeat("A", MaxPermissibleAttachmentSize+1))),
		},
	}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)

	response := golaerror.Error{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(constants.PayloadValidationErrorCode, response.ErrorCode)
	suite.Equal("One or more of the request parameters are missing or invalid", response.ErrorMessage)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldNotThrowBadRequestIfAttachmentSizeIsEqualToAllowedLimit() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{
		{
			FileName: "attachment.pdf",
			Data:     base64.StdEncoding.EncodeToString([]byte(strings.Repeat("A", MaxPermissibleAttachmentSize))),
		},
	}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))
	suite.emailService.EXPECT().Send(suite.context, gomock.Any())

	suite.controller.SendEmail(suite.context)

	suite.NotEqual(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldNotThrowBadRequestForValidEmailRequestWithoutAttachments() {
	request := suite.validEmailRequest()
	request.Attachments = []http_request_response.Attachment{}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))
	suite.emailService.EXPECT().Send(suite.context, gomock.Any())

	suite.controller.SendEmail(suite.context)

	suite.NotEqual(http.StatusBadRequest, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldSendEmailForValidRequest() {
	request := http_request_response.EmailRequest{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: http_request_response.MessageBody{
			Content:  "SGVsbG8h",
			MimeType: "text/html",
		},
		Attachments: []http_request_response.Attachment{{
			FileName: "attachment.pdf",
			Data:     "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
		}},
	}

	suite.emailService.EXPECT().Send(suite.context, models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/html",
			Content:  "Hello!",
		},
		Attachments: []models.Attachment{{
			FileName: "attachment.pdf",
			Data:     []byte("Attachment with some data!"),
		}},
	})

	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	suite.Equal(http.StatusNoContent, suite.recorder.Code)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldRespondWithInternalErrorIfUnableToSendEmail() {
	request := suite.validEmailRequest()

	suite.emailService.EXPECT().Send(suite.context, models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/html",
			Content:  "Hello!",
		},
		Attachments: []models.Attachment{{
			FileName: "attachment.pdf",
			Data:     []byte("Attachment with some data!"),
		}},
	})

	err := &constants.InternalServerError
	suite.emailService.EXPECT().Send(suite.context, gomock.Any()).Return(err)

	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
}

func (suite emailControllerTestSuite) TestSendEmail_ShouldSendEmailWithDefaultMimeType() {
	request := http_request_response.EmailRequest{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: http_request_response.MessageBody{
			Content: "SGVsbG8h",
		},
	}

	suite.emailService.EXPECT().Send(suite.context, gomock.Any())

	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.controller.SendEmail(suite.context)
	suite.Equal(http.StatusNoContent, suite.recorder.Code)
}

func (suite emailControllerTestSuite) validEmailRequest() http_request_response.EmailRequest {
	return http_request_response.EmailRequest{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi",
		Body: http_request_response.MessageBody{
			Content:  "SGVsbG8h",
			MimeType: "text/html",
		},
	}
}

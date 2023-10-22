package controller

import (
	"ccg-api/constants"
	configuration2 "ccg-api/email/configuration"
	"ccg-api/email/http_request_response"
	. "ccg-api/email/service"
	http_util "ccg-api/http-util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
	"net/http"
)

type EmailController interface {
	SendEmail(ctx *gin.Context)
}

type emailController struct {
	config                  configuration2.EmailClientConfig
	service                 EmailService
	httpRequestDeserializer http_util.HttpRequestDeserializer
}

func NewEmailController(service EmailService, config configuration2.EmailClientConfig) EmailController {
	validate := validator.New()

	registerFieldLevelValidator(validate, "validGolaEmailDomain", NewGolaDomainValidator(config.ValidGolaEmailDomain()).validate)
	registerFieldLevelValidator(validate, "validFileExtension", NewFileExtensionValidator(config.UnsupportedAttachmentExtensions()).validate)
	registerFieldLevelValidator(validate, "uniqueAttachments", UniqueAttachmentValidator)
	registerFieldLevelValidator(validate, "notblank", validators.NotBlank)
	registerFieldLevelValidator(validate, "notblankbase64", NewNotBlankBase64ContentValidator().validate)
	registerFieldLevelValidator(validate, "totalAttachmentSizeWithinPermissibleLimit",
		NewTotalAttachmentSizeWithinLimitValidator(config.PermissibleTotalSizeOfAttachments()).Validate)

	return emailController{
		service:                 service,
		httpRequestDeserializer: http_util.NewHttpRequestDeserializer(validate),
		config:                  config,
	}
}

func registerFieldLevelValidator(validate *validator.Validate, tagName string, validateFunction validator.Func) {
	logger := logging.NewLoggerEntry()
	if err := validate.RegisterValidation(tagName, validateFunction); err != nil {
		logger.Fatalf("Failed to register %s validator in email controller", tagName)
	}
}

// SendEmail godoc
// @Tags Email
// @Summary API to send email
// @Description API to send email,
// @Description If IncludeBaseTemplate is true then, header/footer (logos + disclaimer) is included
// @Accept  json
// @Produce  json
// @Param emailRequest body http_request_response.EmailRequest true "Email Request"
// @Success 204
// @Failure 400 {object} golaerror.Error "If From/To/Subject/Body are empty"
// @Failure 500 {object} golaerror.Error ""
// @Router /api/ccg/v1/email/send [post]
func (controller emailController) SendEmail(ctx *gin.Context) {
	// for swagger import
	_ = golaerror.Error{}

	var request http_request_response.EmailRequest
	if bindError := controller.httpRequestDeserializer.ShouldBindJsonBodyIfValid(&request, ctx); bindError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.PayloadValidationError)
		return
	}

	if len(request.Body.MimeType) == 0 {
		request.Body.MimeType = "text/plain"
	}

	email, emailModelError := request.ToEmailModel(ctx)
	if emailModelError != nil {
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	emailSendError := controller.service.Send(ctx, email)
	if emailSendError != nil {
		constants.RespondWithGolaError(ctx, emailSendError)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

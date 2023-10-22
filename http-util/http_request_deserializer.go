package http_util

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/inclusi-blog/gola-utils/logging"
)

type HttpRequestDeserializer interface {
	ShouldBindJsonBodyIfValid(request interface{}, ctx *gin.Context) error
}

type httpRequestDeserializer struct {
	validator *validator.Validate
}

func NewHttpRequestDeserializer(validator *validator.Validate) HttpRequestDeserializer {
	return httpRequestDeserializer{validator: validator}
}

func (deserializer httpRequestDeserializer) ShouldBindJsonBodyIfValid(request interface{}, ctx *gin.Context) error {
	if bindError := ctx.ShouldBindBodyWith(request, binding.JSON); bindError != nil {
		logging.GetLogger(ctx).Error("Failed to bind json request due to error: ", bindError.Error())
		return bindError
	}

	if validationError := deserializer.validator.Struct(request); validationError != nil {
		logging.GetLogger(ctx).Error("Failed to validate request: ", validationError.Error())
		return validationError
	}
	return nil
}

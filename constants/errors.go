package constants

import (
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"net/http"
)

const (
	PayloadValidationErrorCode string = "ERR_CCG_SERVICE_PAYLOAD_INVALID"
	InternalServerErrorCode    string = "ERR_CCG_SERVICE_INTERNAL_SERVER_ERROR"
	CCGServiceFailureCode      string = "ERR_CCG_SERVICE_SERVICE_FAILURE"
)

var (
	CCGServiceFailureError = golaerror.Error{ErrorCode: CCGServiceFailureCode, ErrorMessage: "Failed to communicate with ccg service"}
	PayloadValidationError = golaerror.Error{ErrorCode: PayloadValidationErrorCode, ErrorMessage: "One or more of the request parameters are missing or invalid"}
	InternalServerError    = golaerror.Error{ErrorCode: InternalServerErrorCode, ErrorMessage: "something went wrong"}
)

var ErrorCodeHttpStatusCodeMap = map[string]int{
	PayloadValidationErrorCode: http.StatusBadRequest,
	InternalServerErrorCode:    http.StatusInternalServerError,
	CCGServiceFailureCode:      http.StatusInternalServerError,
}

func GetGolaHttpCode(golaErrCode string) int {
	if httpCode, ok := ErrorCodeHttpStatusCodeMap[golaErrCode]; ok {
		return httpCode
	}
	return http.StatusInternalServerError
}

func RespondWithGolaError(ctx *gin.Context, err error) {
	if golaErr, ok := err.(*golaerror.Error); ok {
		ctx.JSON(GetGolaHttpCode(golaErr.ErrorCode), golaErr)
		return
	}
	ctx.JSON(http.StatusInternalServerError, InternalServerError)
	return
}

package controller

import (
	"ccg-api/email/http_request_response"
	"github.com/go-playground/validator/v10"
)

type TotalAttachmentSizeWithinLimitValidator struct {
	maxTotalAttachmentSize int
}

func NewTotalAttachmentSizeWithinLimitValidator(maxTotalAttachmentSize int) *TotalAttachmentSizeWithinLimitValidator {
	return &TotalAttachmentSizeWithinLimitValidator{maxTotalAttachmentSize: maxTotalAttachmentSize}
}

func (validator TotalAttachmentSizeWithinLimitValidator) Validate(fieldLevel validator.FieldLevel) bool {
	attachments := fieldLevel.Field().Interface().([]http_request_response.Attachment)
	totalSize := 0
	for _, attachment := range attachments {
		if data, err := attachment.GetDecodedData(); err == nil {
			totalSize = totalSize + len(data)
		}
	}
	return totalSize <= validator.maxTotalAttachmentSize
}

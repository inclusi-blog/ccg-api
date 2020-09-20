package controller

import (
	httprequestresponse "ccg-api/email/http_request_response"
	"github.com/go-playground/validator/v10"
)

func UniqueAttachmentValidator(fieldLevel validator.FieldLevel) bool {
	attachments := fieldLevel.Field().Interface().([]httprequestresponse.Attachment)
	uniqueFiles := map[string]bool{}
	for _, attachment := range attachments {
		if _, duplicateFileName := uniqueFiles[attachment.FileName]; duplicateFileName {
			return false
		}
		uniqueFiles[attachment.FileName] = true
	}
	return true
}

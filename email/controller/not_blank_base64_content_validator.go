package controller

import (
	"encoding/base64"
	"github.com/go-playground/validator/v10"
	"strings"
)

type NotBlankBase64ContentValidator struct {
}

func NewNotBlankBase64ContentValidator() *NotBlankBase64ContentValidator {
	return &NotBlankBase64ContentValidator{}
}

func (validator NotBlankBase64ContentValidator) validate(fieldLevel validator.FieldLevel) bool {
	value := fieldLevel.Field().String()
	if decodedBytes, err := base64.StdEncoding.DecodeString(value); err == nil {
		return len(strings.TrimSpace(string(decodedBytes))) > 0
	}
	return false
}

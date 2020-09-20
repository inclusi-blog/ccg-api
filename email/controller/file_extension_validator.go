package controller

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

type FileExtensionValidator struct {
	unsupportedExtensions []string
}

func NewFileExtensionValidator(unsupportedExtensions []string) *FileExtensionValidator {
	return &FileExtensionValidator{unsupportedExtensions: unsupportedExtensions}
}

func (fileExtensionValidator FileExtensionValidator) validate(fieldLevel validator.FieldLevel) bool {
	fileName := strings.TrimSpace(fieldLevel.Field().String())
	fileExtensionTypeIndex := strings.LastIndex(fileName, ".")

	if fileExtensionTypeIndex <= 0 || fileExtensionTypeIndex == len(fileName)-1 {
		return false
	}

	for _, extension := range fileExtensionValidator.unsupportedExtensions {
		if strings.HasSuffix(fileName, "."+extension) {
			return false
		}
	}
	return true
}

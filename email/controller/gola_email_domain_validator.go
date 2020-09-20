package controller

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

type GolaDomainValidator struct {
	validGolaDomain []string
}

func NewGolaDomainValidator(validGolaDomains []string) *GolaDomainValidator {
	return &GolaDomainValidator{validGolaDomain: validGolaDomains}
}

func (golaDomainValidator GolaDomainValidator) validate(fieldLevel validator.FieldLevel) bool {
	domainInRequest := strings.ToLower(fieldLevel.Field().String())
	for _, domain := range golaDomainValidator.validGolaDomain {
		domain = strings.ToLower(domain)
		if strings.HasSuffix(domainInRequest, "@"+domain) {
			return true
		}
	}
	return false
}

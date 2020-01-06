package api

import (
	"github.com/asaskevich/govalidator"

	"github.com/blackplayerten/IdealVisual_backend/account"
)

const (
	WrongEmail = "not_email"

	WrongLen = "wrong_len"
)

func validateAll(info *account.FullAccount) []*FieldError {
	validationErrors := make([]*FieldError, 0, 3)

	if validationError := validateEmail(info.Email); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError := validatePassword(info.Password); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError := validateUsername(info.Username); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	return validationErrors
}

func validateEmail(e string) *FieldError {
	if isValid := govalidator.IsEmail(e); !isValid {
		return &FieldError{
			Field:   "email",
			Reasons: []string{WrongEmail},
		}
	}

	return nil
}

func validatePassword(p string) *FieldError {
	if isValid := govalidator.StringLength(p, "8", "64"); !isValid {
		return &FieldError{
			Field:   "password",
			Reasons: []string{WrongLen},
		}
	}

	return nil
}

func validateUsername(u string) *FieldError {
	if isValid := govalidator.StringLength(u, "4", "32"); !isValid {
		return &FieldError{
			Field:   "username",
			Reasons: []string{WrongLen},
		}
	}

	return nil
}

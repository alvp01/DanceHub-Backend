// internal/validator/password.go
package validator

import (
	"errors"
	"unicode"
)

var ErrPasswordTooShort = errors.New("el password debe tener al menos 12 caracteres")
var ErrPasswordNeedUppercase = errors.New("el password debe tener al menos 2 letras mayúsculas")
var ErrPasswordNeedDigits = errors.New("el password debe tener al menos 3 números")
var ErrPasswordNeedSpecial = errors.New("el password debe tener al menos 2 caracteres especiales (*-!_.^)")

var allowedSpecial = map[rune]bool{
	'*': true, '-': true, '!': true,
	'_': true, '.': true, '^': true,
}

type PasswordValidationError struct {
	Errors []string `json:"errors"`
}

func (e *PasswordValidationError) Error() string {
	return "password inválido"
}

func ValidatePassword(password string) error {
	var (
		upperCount   int
		digitCount   int
		specialCount int
	)

	validationErr := &PasswordValidationError{}

	if len(password) < 12 {
		validationErr.Errors = append(validationErr.Errors, ErrPasswordTooShort.Error())
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upperCount++
		case unicode.IsDigit(char):
			digitCount++
		case allowedSpecial[char]:
			specialCount++
		}
	}

	if upperCount < 2 {
		validationErr.Errors = append(validationErr.Errors, ErrPasswordNeedUppercase.Error())
	}
	if digitCount < 3 {
		validationErr.Errors = append(validationErr.Errors, ErrPasswordNeedDigits.Error())
	}
	if specialCount < 2 {
		validationErr.Errors = append(validationErr.Errors, ErrPasswordNeedSpecial.Error())
	}

	if len(validationErr.Errors) > 0 {
		return validationErr
	}

	return nil
}

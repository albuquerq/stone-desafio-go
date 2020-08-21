package validator

import (
	valid "github.com/go-ozzo/ozzo-validation/v4"
	match "github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/access"
)

// ValidateAccessCredential validates access credentials data.
func ValidateAccessCredential(cred access.Credential) error {
	return valid.ValidateStruct(&cred,
		valid.Field(&cred.CPF, valid.Required, valid.Length(11, 11), match.Digit),
		valid.Field(&cred.Secret, valid.Required),
	)
}

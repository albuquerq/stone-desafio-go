package validator

import (
	"strings"

	valid "github.com/go-ozzo/ozzo-validation/v4"
	match "github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type acCreationValidator struct {
	log *logrus.Entry
}

// NewAccountCreation validates whether a account can be created.
func NewAccountCreation() account.ValidationStrategy {
	return &acCreationValidator{
		log: common.Logger().WithField("source", "accountCreationValidator"),
	}
}

func (acv *acCreationValidator) Validate(ac account.Account) error {
	ac.Name = strings.TrimSpace(ac.Name) // It will not accept strings of blanks.

	err := valid.ValidateStruct(&ac,
		valid.Field(&ac.ID, valid.Required, match.UUID),
		valid.Field(&ac.Name, valid.Required),
		valid.Field(&ac.CPF, valid.Required, valid.Length(11, 11), match.Digit),
		valid.Field(&ac.Balance, valid.Min(0)), // Will not accept accounts with negative balance.
		valid.Field(&ac.Secret, valid.Required),
	)

	if err != nil {
		acv.log.WithError(err).Error("data validation err")
		return err
	}

	return nil
}

type acBalanceUpdateValidator struct {
	log *logrus.Entry
}

// NewAccountBalanceUpdate validates data for updating an account balance.
func NewAccountBalanceUpdate() account.ValidationStrategy {
	return &acBalanceUpdateValidator{
		log: common.Logger().WithField("source", "acBalanceUpdateValidator"),
	}
}

func (abuv *acBalanceUpdateValidator) Validate(ac account.Account) error {
	err := valid.ValidateStruct(&ac,
		valid.Field(&ac.Balance, valid.Min(0)), //  Will not accept accounts with negative balance.
	)
	if err != nil {
		abuv.log.WithError(err).Error("data validation err")
	}
	return nil
}

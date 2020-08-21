package validator

import (
	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type transferCreationValidator struct {
	log *logrus.Entry
}

// NewTransferCreation validates whether a transfer can be created.
func NewTransferCreation() transfer.ValidationStrategy {
	return &transferCreationValidator{
		log: common.Logger().WithFields(logrus.Fields{
			"source": "TransferCreationValidator",
			"op":     "Validate",
		}),
	}
}

func (tv transferCreationValidator) Validate(tr transfer.Transfer) error {
	log := tv.log.WithFields(logrus.Fields{
		"accountOringinID":     tr.AccountOriginID,
		"accountDestinationID": tr.AccountDestinationID,
		"amount":               tr.Amount,
	})

	if tr.ID == "" {
		return errors.ErrNoHasUniqueIdentity
	}

	if tr.AccountOriginID == "" {
		return errors.ErrTransferMissingData
	}

	if tr.AccountDestinationID == "" {
		return errors.ErrTransferMissingData
	}

	if tr.AccountOriginID == tr.AccountOriginID {
		err := errors.ErrTransferBetweenSameAccount
		log.WithError(err)
		return err
	}

	if tr.Amount == 0 {
		err := errors.ErrTransferMissingAmount
		log.WithError(err)
		return err
	}

	return nil
}

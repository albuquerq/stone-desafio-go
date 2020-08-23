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
		err := errors.ErrNoHasUniqueIdentity
		log.WithError(err).Error("missing transfer uuid")
		return err
	}

	if tr.AccountOriginID == "" {
		err := errors.ErrTransferMissingData
		log.WithError(err).Error("missing origin account id")
	}

	if tr.AccountDestinationID == "" {
		err := errors.ErrTransferMissingData
		log.WithError(err).Error("missing destination account id")
		return err
	}

	if tr.AccountOriginID == tr.AccountDestinationID {
		err := errors.ErrTransferBetweenSameAccount
		log.WithError(err).Error("transfer betwen some account")
		return err
	}

	if tr.Amount == 0 {
		err := errors.ErrTransferMissingAmount
		log.WithError(err).Error("amount not defined")
		return err
	}

	return nil
}

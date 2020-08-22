package application

import (
	goerrors "errors"

	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type transferService struct {
	repoRegistry      domain.RepositoryRegistry
	trCreateValidator transfer.ValidationStrategy
	log               *logrus.Entry
}

// NewTransferService return a new transfer.Service.
func NewTransferService(
	repoRegistry domain.RepositoryRegistry,
	trCreateValid transfer.ValidationStrategy,
) transfer.Service {

	return &transferService{
		repoRegistry:      repoRegistry,
		trCreateValidator: trCreateValid,
		log:               common.Logger().WithField("source", "transferService"),
	}
}

func (ts *transferService) checkDependencies() {
	if ts.repoRegistry == nil {
		ts.log.Fatal("RepositoryRegistry is not defined")
	}

	if ts.trCreateValidator == nil {
		ts.log.Fatal("TransferCreateValidator is not defined")
	}
}

func (ts *transferService) Transfer(fromID string, toID string, amount int) (tr transfer.Transfer, err error) {
	ts.checkDependencies()

	log := ts.log.WithFields(logrus.Fields{
		"op":                     "Transfer",
		"account_origin_id":      fromID,
		"account_destination_id": toID,
		"amount":                 amount,
	})

	acRepository := ts.repoRegistry.AccountRepository()

	acFrom, err := acRepository.GetByID(fromID)
	if err != nil {
		log.WithError(err).Error("origin account not found")
		if goerrors.Is(err, errors.ErrAccountNotFound) {
			err = errors.ErrTransferAccountOriginNotFound
		}
		return
	}

	acTo, err := acRepository.GetByID(toID)
	if err != nil {
		log.WithError(err).Error("destination account not found")
		if goerrors.Is(err, errors.ErrAccountNotFound) {
			err = errors.ErrTransferAccountDestinationNotFound
		}
		return
	}

	if acFrom.Balance < amount {
		err = errors.ErrTransferInsufficientBalance
		ts.log.WithError(err).Error("oringin balance is insufficient for transfer")
		return
	}

	acFrom.Balance -= amount
	acTo.Balance += amount

	tx := ts.repoRegistry.Tx()

	err = acRepository.WithTx(tx).UpdateBalance(acFrom)
	if err != nil {
		log.WithError(err).Error("failed to update the source account balance")
		tx.Rollback()
		return
	}

	err = acRepository.WithTx(tx).UpdateBalance(acTo)
	if err != nil {
		log.WithError(err).Error("failed to update target account balance")
		tx.Rollback()
		return
	}

	trRepository := ts.repoRegistry.TransferRepository()

	tr = transfer.Transfer{
		ID:                   trRepository.GenerateIdentifier(),
		AccountOriginID:      acFrom.ID,
		AccountDestinationID: acTo.ID,
		Amount:               amount,
	}

	err = trRepository.WithTx(tx).Store(&tr)
	if err != nil {
		log.WithError(err).Error("failed to store the transfer")
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.WithError(err).Error("failed to commit data")
		return
	}

	return
}

func (ts *transferService) ListTransfersFromAccount(accountID string) ([]transfer.Transfer, error) {
	log := ts.log.WithField("op", "ListTransfersFromAccount").
		WithField("account_id", accountID)

	transfers, err := ts.repoRegistry.TransferRepository().ListByAccountID(accountID)
	if err != nil {
		log.WithError(err).Error("an error occurred while listing account transfers")
	}

	return transfers, nil
}

package mem

import (
	"sync"
	"time"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/utils"

	"github.com/sirupsen/logrus"
)

type memTransferRepo struct {
	log       *logrus.Entry
	transfers []transfer.Transfer
	mux       sync.Mutex
}

// NewTransferRepository return a new in-memory transfer.Repository.
func NewTransferRepository(log *logrus.Logger) transfer.Repository {
	return &memTransferRepo{
		log:       log.WithField("source", "memTransferRepo"),
		transfers: []transfer.Transfer{},
	}
}

func (mtr *memTransferRepo) Store(tr *transfer.Transfer) error {

	mtr.mux.Lock()
	defer mtr.mux.Unlock()

	tr.CreatedAt = time.Now()

	mtr.transfers = append(mtr.transfers, *tr)

	return nil
}

func (mtr *memTransferRepo) GetById(transferID string) (tr transfer.Transfer, err error) {

	index := mtr.indexOf(transferID)
	if index < 0 {
		err = errors.ErrTransferNotFound
		mtr.log.WithError(err).WithField("transferID", transferID)
		return
	}

	mtr.mux.Lock()
	defer mtr.mux.Unlock()

	tr = mtr.transfers[index]

	return
}

func (mtr *memTransferRepo) indexOf(transferID string) int {
	mtr.mux.Lock()
	defer mtr.mux.Unlock()

	for i, tr := range mtr.transfers {
		if tr.ID == transferID {
			return i
		}
	}
	return -1
}

func (mtr *memTransferRepo) ListByAccountID(accountID string) ([]transfer.Transfer, error) {
	transfers := []transfer.Transfer{}

	mtr.mux.Lock()
	defer mtr.mux.Unlock()

	for _, tr := range mtr.transfers {
		if accountID == tr.AccountOriginID || accountID == tr.AccountDestinationID {
			localTr := tr // make local copy
			transfers = append(transfers, localTr)
		}
	}

	return transfers, nil
}

func (mtr *memTransferRepo) GenerateIndetifier() string {
	return utils.GenUUID()
}
package mem

import (
	"database/sql/driver"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type memAccountRepo struct {
	accounts []account.Account
	mux      sync.RWMutex
	log      *logrus.Entry
}

// NewAccoutRepository returns an in-memory repository for account.Account.
func NewAccoutRepository(logger *logrus.Logger) account.Repository {
	return &memAccountRepo{
		log:      logger.WithField("source", "memAccountRepo"),
		accounts: []account.Account{},
	}
}

func (mar *memAccountRepo) Store(ac *account.Account) error {
	log := mar.log.WithField("op", "Store")

	if ac.ID == "" {
		err := errors.ErrNoHasUniqueIdentity
		log.WithError(err).Error("account ID not defined")
		return err
	}
	//mar.log.Infof("Storing a account %v", ac.ID)

	ac.CreatedAt = time.Now()

	mar.mux.Lock()
	mar.accounts = append(mar.accounts, *ac)
	mar.mux.Unlock()

	log.Infof("account %s successfully stored", ac.ID)

	return nil
}

func (mar *memAccountRepo) UpdateBalance(ac account.Account) error {
	log := mar.log.WithField("op", "UpdateBalance")

	if ac.ID == "" {
		err := errors.ErrNoHasUniqueIdentity
		log.WithError(err).WithField("accountID", ac.ID)
		return err
	}

	index := mar.indexOf(ac.ID)
	if index < 0 {
		err := errors.ErrAccountNotFound
		log.WithError(err).WithField("accountID", ac.ID)
		return err
	}

	mar.mux.Lock()
	defer mar.mux.Unlock()
	mar.accounts[index].Balance = ac.Balance

	return nil
}

func (mar *memAccountRepo) indexOf(accountID string) int {
	mar.mux.RLock()
	defer mar.mux.RUnlock()

	for i, ac := range mar.accounts {
		if ac.ID == accountID {
			return i
		}
	}
	return -1
}

func (mar *memAccountRepo) indexOfCPF(cpf string) int {
	mar.mux.RLock()
	defer mar.mux.RUnlock()

	for i, ac := range mar.accounts {
		if cpf == ac.CPF {
			return i
		}
	}
	return -1
}

func (mar *memAccountRepo) GetByID(acID string) (ac account.Account, err error) {
	log := mar.log.WithField("op", "GetByID")

	if acID == "" {
		err = errors.ErrNoHasUniqueIdentity
		log.WithError(err).Error("accountID", acID)
		return
	}

	index := mar.indexOf(acID)
	if index < 0 {
		err = errors.ErrAccountNotFound
		log.WithError(err).WithField("accountID", acID)
		return
	}
	mar.mux.RLock()
	ac = mar.accounts[index]
	mar.mux.RUnlock()

	return ac, nil
}

func (mar *memAccountRepo) GetByCPF(cpf string) (ac account.Account, err error) {
	log := mar.log.WithField("op", "GetByCPF")

	if cpf == "" {
		err = errors.ErrNoHasUniqueIdentity
		log.WithError(err).WithField("cpf", cpf)
		return
	}

	index := mar.indexOfCPF(cpf)
	if index < 0 {
		err = errors.ErrAccountNotFound
		log.WithError(err).WithField("cpf", cpf).Error("account with cpf not found")
		return
	}

	mar.mux.RLock()
	ac = mar.accounts[index]
	mar.mux.RUnlock()

	return ac, err
}

func (mar *memAccountRepo) ListAll() ([]account.Account, error) {
	mar.mux.RLock()
	defer mar.mux.RUnlock()

	accounts := make([]account.Account, len(mar.accounts))

	copy(accounts, mar.accounts)

	return accounts, nil
}

func (mar *memAccountRepo) GenerateIdentifier() string {
	return common.GenUUID()
}

func (mar *memAccountRepo) WithTx(tx driver.Tx) account.Repository {
	mar.log.WithField("op", "WithTx").Debug("memory transactions are not applicable")
	return mar
}

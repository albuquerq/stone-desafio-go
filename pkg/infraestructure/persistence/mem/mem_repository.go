package mem

import (
	"sync"
	"time"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/utils"
	"github.com/sirupsen/logrus"
)

type memAccountRepository struct {
	accounts []account.Account
	mux      sync.Mutex
	log      *logrus.Entry
}

// NewAccoutRepository returns an in-memory repository for account.Account.
func NewAccoutRepository(logger *logrus.Logger) account.Repository {
	return &memAccountRepository{
		log:      logger.WithField("source", "memAccountRepository"),
		accounts: []account.Account{},
	}
}

func (mar *memAccountRepository) Store(ac *account.Account) error {
	if ac.ID == "" {
		err := errors.ErrNoHasUniqueIdentity
		mar.log.WithError(err).Error("account ID not defined")
		return err
	}
	//mar.log.Infof("Storing a account %v", ac.ID)

	ac.CreatedAt = time.Now()

	mar.mux.Lock()
	defer mar.mux.Unlock()

	mar.accounts = append(mar.accounts, *ac)

	mar.log.Infof("account %s successfully stored", ac.ID)

	return nil
}

func (mar *memAccountRepository) UpdateBalance(ac account.Account) error {
	if ac.ID == "" {
		err := errors.ErrNoHasUniqueIdentity
		mar.log.WithError(err).WithField("accountID", ac.ID)
		return err
	}

	index := mar.indexOf(ac.ID)
	if index < 0 {
		err := errors.ErrAccountNotFound
		mar.log.WithError(err).WithField("accountID", ac.ID)
		return err
	}

	mar.mux.Lock()
	defer mar.mux.Unlock()

	mar.accounts[index].Balance = ac.Balance

	return nil
}

func (mar *memAccountRepository) indexOf(accountID string) int {
	mar.mux.Lock()
	defer mar.mux.Unlock()

	for i, ac := range mar.accounts {
		if ac.ID == accountID {
			return i
		}
	}
	return -1
}

func (mar *memAccountRepository) GetByID(acID string) (account.Account, error) {
	if acID == "" {
		err := errors.ErrNoHasUniqueIdentity
		mar.log.WithError(err).Error("accountID", acID)
		return account.Account{}, err
	}

	index := mar.indexOf(acID)
	if index < 0 {
		err := errors.ErrAccountNotFound
		mar.log.WithError(err).WithField("accountID", acID)
		return account.Account{}, err
	}
	mar.mux.Lock()
	defer mar.mux.Unlock()

	ac := mar.accounts[index]

	return ac, nil
}

func (mar *memAccountRepository) ListAll() ([]account.Account, error) {
	mar.mux.Lock()
	defer mar.mux.Unlock()

	accounts := make([]account.Account, len(mar.accounts))

	copy(accounts, mar.accounts)

	return accounts, nil
}

func (mar *memAccountRepository) GenerateIdentifier() string {
	return utils.GenUUID()
}

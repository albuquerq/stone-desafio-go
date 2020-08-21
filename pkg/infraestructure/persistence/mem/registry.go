package mem

import (
	"database/sql/driver"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

// NewRepositoryRegistry reutns a new MemRepositoryRegistry.
func NewRepositoryRegistry() domain.RepositoryRegistry {
	return &memRepositoryRegistry{log: common.Logger()}
}

type memRepositoryRegistry struct {
	log *logrus.Logger
}

var (
	acRepo account.Repository
	acOnce sync.Once

	trRepo transfer.Repository
	trOnce sync.Once
)

func (rr *memRepositoryRegistry) AccountRepository() account.Repository {
	acOnce.Do(func() { // Lazy singleton.
		acRepo = NewAccoutRepository()
	})
	return acRepo
}

func (rr *memRepositoryRegistry) TransferRepository() transfer.Repository {
	trOnce.Do(func() { // Lazy singleton.
		trRepo = NewTransferRepository()
	})
	return trRepo
}

func (rr *memRepositoryRegistry) Tx() driver.Tx {
	return &memTx{}
}

package application

import (
	"database/sql/driver"

	"github.com/albuquerq/stone-desafio-go/pkg/domain"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/access"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/validator"
)

type registry struct {
	accountService  account.Service
	transferService transfer.Service
	accessService   access.Service

	repoRegistry domain.RepositoryRegistry
}

// NewRegistry returns a new application registry.
func NewRegistry(
	repoRegistry domain.RepositoryRegistry,
) domain.Registry {

	r := registry{
		repoRegistry: repoRegistry,
		accountService: NewAccountService(
			repoRegistry,
			validator.NewAccountCreation(),
			validator.NewAccountBalanceUpdate(),
		),
		transferService: NewTransferService(
			repoRegistry,
			validator.NewTransferCreation(),
		),
		accessService: NewAccessService(
			repoRegistry,
		),
	}
	r.checkDependencies()

	return &r
}

func (r *registry) checkDependencies() {
	if r.accessService == nil {
		panic("access service not defined")
	}
	if r.accountService == nil {
		panic("account service not defined")
	}
	if r.transferService == nil {
		panic("transfer service not defined")
	}
	if r.repoRegistry == nil {
		panic("repository registry not defined")
	}
}

func (r *registry) AccessService() access.Service {
	return r.accessService
}

func (r *registry) AccountService() account.Service {
	return r.accountService
}

func (r *registry) TransferService() transfer.Service {
	return r.transferService
}

func (r *registry) AccountRepository() account.Repository {
	return r.repoRegistry.AccountRepository()
}

func (r *registry) TransferRepository() transfer.Repository {
	return r.repoRegistry.TransferRepository()
}

func (r *registry) Tx() driver.Tx {
	return r.repoRegistry.Tx()
}

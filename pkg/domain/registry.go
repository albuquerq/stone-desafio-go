package domain

import (
	"database/sql/driver"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/access"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
)

// RepositoryRegistry group repository factories.
type RepositoryRegistry interface {
	AccountRepository() account.Repository
	TransferRepository() transfer.Repository
	Tx() driver.Tx
}

// ServiceRegistry group service factories.
type ServiceRegistry interface {
	AccessService() access.Service
	AccountService() account.Service
	TransferService() transfer.Service
}

type Registry interface {
	RepositoryRegistry
	ServiceRegistry
}

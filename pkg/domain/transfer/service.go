package transfer

import "github.com/albuquerq/stone-desafio-go/pkg/domain/account"

// Service encapsulates domain operations over the transfer model as a service.
type Service interface {
	Transfer(from account.Account, to account.Account) (Transfer, error)
	ListTransfersFromAccount(accountID string) ([]Transfer, error)
}

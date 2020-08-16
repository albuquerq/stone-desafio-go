package access

import "github.com/albuquerq/stone-desafio-go/pkg/domain/account"

// Service responsible for authentication.
type Service interface {
	Authenticate(Credential) (account.Account, error)
}

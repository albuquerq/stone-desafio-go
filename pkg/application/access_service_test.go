package application

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/access"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

var (
	accessServiceVar access.Service = NewAccessService(repoRegistry) // repoRegistry defined in "account_service_test.go" file.
)

func TestAcessService(t *testing.T) {

	ac := account.Account{
		ID:      common.GenUUID(),
		Name:    "Jhon Doe Acess",
		Balance: 500050,
		CPF:     "76312890312",
		Secret:  "secret pass 1",
	}

	// Create account for test.
	err := acService.CreateAccount(&ac)
	assert.NoError(t, err)

	acDescription, err := accessServiceVar.Authenticate(access.Credential{
		CPF:    "76312890312",
		Secret: "err secret",
	})
	assert.Error(t, err)

	acDescription, err = accessServiceVar.Authenticate(access.Credential{
		CPF:    "76312890312",
		Secret: "secret pass 1",
	})
	assert.NoError(t, err)

	assert.Equal(t, ac.ID, acDescription.AccountID)
	assert.Equal(t, "76312890312", acDescription.CPF)
	assert.Equal(t, ac.Name, acDescription.Name)

}

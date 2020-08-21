package application

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/persistence/mem"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/validator"
)

//  Account service test using in-memory account.Repository implementation and validators using ozzo-validator.

var (
	repoRegistry                  = mem.NewRepositoryRegistry()
	accountCreationValidator      = validator.NewAccountCreation()
	accountBalanceUpdateValidator = validator.NewAccountBalanceUpdate()

	acService = NewAccountService(
		repoRegistry,
		accountCreationValidator,
		accountBalanceUpdateValidator,
	)
)

func TestAccountService(t *testing.T) {

	ac := account.Account{
		ID:      repoRegistry.AccountRepository().GenerateIdentifier(),
		CPF:     "00000000000",
		Name:    "Jon Due",
		Balance: 100,
		Secret:  "some secret",
	}

	t.Run("CreateAccount", func(t *testing.T) {
		err := acService.CreateAccount(&ac)
		assert.NoError(t, err)
	})

	t.Run("UpdateBalance", func(t *testing.T) {
		oldBalance := ac.Balance

		ac.Balance = 10000

		err := acService.UpdateBalance(&ac)
		assert.NoError(t, err)

		var balanceUpdated account.BalanceValue

		t.Run("AccountBalance", func(t *testing.T) {
			balanceUpdated, err = acService.AccountBalance(ac.ID)
			assert.NoError(t, err)
		})

		assert.NotEqual(t, oldBalance, balanceUpdated.Balance)

		assert.Equal(t, ac.Balance, balanceUpdated.Balance)
	})

	t.Run("ListAccounts", func(t *testing.T) {
		accounts, err := acService.ListAccounts()
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(accounts))

		for _, ac := range accounts {
			t.Logf("Account %s - Banlance R$ %.2f", ac.ID, float32(ac.Balance/100))
		}
	})
}

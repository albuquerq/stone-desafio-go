package application

import (
	goerrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/persistence/mem"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/validator"
)

var (
	trService = NewTransferService(
		mem.NewRepositoryRegistry(),
		validator.NewTransferCreation(),
	)
)

func TestTransferService(t *testing.T) {

	ac1 := account.Account{
		Name:    "John Doe",
		CPF:     "00000000000",
		Balance: 100000,
		Secret:  "123",
	}

	ac2 := account.Account{
		Name:    "Jane Doe",
		CPF:     "11111111111",
		Balance: 50000,
		Secret:  "1234",
	}

	err := acService.CreateAccount(&ac1)
	assert.NoError(t, err)

	err = acService.CreateAccount(&ac2)
	assert.NoError(t, err)

	var transfers []transfer.Transfer

	t.Run("checks that c1 no have transfers", func(t *testing.T) {

		transfers, err = trService.ListTransfersFromAccount(ac1.ID)
		assert.NoError(t, err)

	})

	t.Run("transfers R$ 250 from ac1 to ac2", func(t *testing.T) {

		assert.Equal(t, 0, len(transfers))

		tr, err := trService.Transfer(ac1.ID, ac2.ID, 25000)
		assert.NoError(t, err)

		assert.Equal(t, tr.AccountOriginID, ac1.ID)
		assert.Equal(t, tr.AccountDestinationID, ac2.ID)
		assert.Equal(t, 25000, tr.Amount)
	})

	t.Run("retrieves the values of c1 and c2 and checks whether their balances have been changed correctly", func(t *testing.T) {

		ac1, err = acService.GetAccount(ac1.ID)
		assert.NoError(t, err)

		ac2, err = acService.GetAccount(ac2.ID)
		assert.NoError(t, err)

		assert.Equal(t, 75000, ac1.Balance)
		assert.Equal(t, 75000, ac2.Balance)

	})

	transfers, err = trService.ListTransfersFromAccount(ac1.ID)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(transfers))

	assert.Equal(t, 25000, transfers[0].Amount)

	tr2, err := trService.Transfer(ac1.ID, ac2.ID, 200000)
	assert.Error(t, err)

	assert.Equal(t, true, goerrors.Is(err, errors.ErrTransferInsufficientBalance))

	assert.Equal(t, tr2.ID, "")
}

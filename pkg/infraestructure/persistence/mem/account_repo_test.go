package mem

import (
	"math/rand"
	"testing"
	"time"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"

	"github.com/stretchr/testify/assert"
)

type testCaseSingle struct {
	Case          string
	In            account.Account
	ExpectedValue account.Account
	ExpectedError error
}

type testCaseGrouped struct {
	Case     string
	Expected []account.Account
}

var passStoreTestCases = []testCaseSingle{
	{
		Case: "Valid account",
		In: account.Account{
			ID:      common.GenUUID(),
			Name:    "Jon Due",
			Balance: 0,
			CPF:     "00000000003",
			Secret:  "a secret",
		},
		ExpectedError: nil,
	},
	{
		Case: "With the CreatedAt field defined, it ignores and replaces by insertion time",
		In: account.Account{
			ID:        common.GenUUID(),
			Name:      "Jon Due 2",
			Balance:   0,
			CPF:       "00000004000",
			Secret:    "a secret",
			CreatedAt: time.Now().Add(2 * time.Hour),
		},
		ExpectedError: nil,
	},
}

var failsStoreTestCases = []testCaseSingle{
	{
		Case: "ID not defined",
		In: account.Account{
			Name:    "Jon Due 3",
			CPF:     "12345678910",
			Balance: 1000,
			Secret:  "a secret",
		},
		ExpectedError: errors.ErrNoHasUniqueIdentity,
	},
	{
		Case: "Not stored account",
		In: account.Account{
			ID:   common.GenUUID(),
			Name: "Jon Due 5",
		},
		ExpectedError: errors.ErrAccountNotFound,
	},
}

var (
	memAccountRepository = NewAccoutRepository()
)

// Black-box test.
func TestAccountRepository_Store_And_GetByID(t *testing.T) {

	t.Run("pass", func(t *testing.T) {

		t.Run("account.Repository.Store test", func(t *testing.T) {
			for i := range passStoreTestCases {
				t.Log(passStoreTestCases[i].Case)
				assert.NoError(t, memAccountRepository.Store(&passStoreTestCases[i].In))
			}
		})

		t.Run("account.Repository.GetByID and GetByCPF test", func(t *testing.T) {
			for _, tc := range passStoreTestCases {
				t.Log(tc.Case)
				ac, err := memAccountRepository.GetByID(tc.In.ID)

				assert.NoError(t, err)

				assert.Equal(t, tc.In, ac)
			}

			for _, tc := range passStoreTestCases {
				t.Log(tc.Case)
				ac, err := memAccountRepository.GetByCPF(tc.In.CPF)

				assert.NoError(t, err)
				assert.Equal(t, tc.In.CPF, ac.CPF)
			}
		})
	})

	t.Run("fail", func(t *testing.T) {

		t.Run("account.Repository.Store", func(t *testing.T) {

			for i, tc := range failsStoreTestCases {
				if i == len(failsStoreTestCases)-1 { // Not store the last item.
					break
				}
				t.Log(tc.Case)
				err := memAccountRepository.Store(&tc.In)
				assert.Error(t, err)
			}
		})

		t.Run("account.Respository.GetByID", func(t *testing.T) {

			// Test not stored account.
			tc := failsStoreTestCases[len(failsStoreTestCases)-1]
			t.Log(tc.Case)
			_, err := memAccountRepository.GetByID(tc.In.ID)
			if assert.Error(t, err) {
				assert.Equal(t, errors.ErrAccountNotFound, err)
			}

			// Test blank ID get.
			_, err = memAccountRepository.GetByID("")
			if assert.Error(t, err) {
				assert.Equal(t, errors.ErrNoHasUniqueIdentity, err)
			}
		})

	})

}

func TestAccountRespository_UpdateBalance(t *testing.T) {
	// Warnig: Do not perform this test individually. It depends on the tests above.

	accounts, err := memAccountRepository.ListAll()

	assert.NoError(t, err)

	assert.True(t, len(accounts) > 0)

	ac := accounts[0]

	oldBalance := ac.Balance

	ac.Balance = rand.Int()

	err = memAccountRepository.UpdateBalance(ac)

	assert.NoError(t, err)

	acNew, err := memAccountRepository.GetByID(ac.ID)
	assert.NoError(t, err)

	assert.Equal(t, acNew.Balance, ac.Balance)

	assert.NotEqual(t, oldBalance, ac.Balance)
}

func TestRepository_ListAll(t *testing.T) {
	accounts, err := memAccountRepository.ListAll()

	assert.NoError(t, err)

	for _, ac := range accounts {
		t.Logf("Account %s -> %d BLR cents", ac.ID, ac.Balance)
	}

	assert.Equal(t, len(accounts), len(failsStoreTestCases))
}

func TestAccountRepository_GenerateIdentifier(t *testing.T) {
	newID := memAccountRepository.GenerateIdentifier()

	assert.NotEqual(t, "", newID)
}

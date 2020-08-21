package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type transferRepoTestCase struct {
	Description   string
	In            transfer.Transfer
	ExpectedError error
}

var (
	passCases = []transferRepoTestCase{
		{
			Description: "Valid transfer",
			In: transfer.Transfer{
				ID:                   common.GenUUID(),
				AccountOriginID:      common.GenUUID(),
				AccountDestinationID: common.GenUUID(),
				Amount:               100,
			},
		},
	}
)

var (
	transferMemRepository = NewTransferRepository(logrus.New())
)

func TestTransferRepository(t *testing.T) {

	t.Run("TransferMemRepository.Store and GetByID", func(t *testing.T) {
		for _, tc := range passCases {
			t.Log(tc.Description)
			initCreatedDate := tc.In.CreatedAt
			err := transferMemRepository.Store(&tc.In)

			assert.NoError(t, err)

			assert.NotEqual(t, initCreatedDate, tc.In.CreatedAt)

			tr2, err := transferMemRepository.GetByID(tc.In.ID)

			assert.NoError(t, err)

			assert.Equal(t, tc.In, tr2)

		}
	})

	t.Run("Transfer not found", func(t *testing.T) {
		_, err := transferMemRepository.GetByID(common.GenUUID())
		assert.Error(t, err)
	})

	t.Run("List tansfers", func(t *testing.T) {
		transfers, err := transferMemRepository.ListByAccountID(passCases[0].In.AccountOriginID)
		assert.NoError(t, err)

		assert.NotEqual(t, 0, len(transfers))
	})

	t.Run("Test generate id", func(t *testing.T) {
		id := transferMemRepository.GenerateIndetifier()
		assert.NotEmpty(t, id)
	})

}

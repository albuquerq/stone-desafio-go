package psql

import (
	"database/sql"
	"database/sql/driver"

	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type pgTransferRepo struct {
	db  *sql.DB
	log *logrus.Entry
	tx  driver.Tx
}

// NewTransferRepository returns a transfer repository implementation for PostgreSQL.
func NewTransferRepository(db *sql.DB) transfer.Repository {
	return &pgTransferRepo{
		log: common.Logger().WithField("source", "pgTrasnferRepo"),
		db:  db,
	}
}

func (pgtr *pgTransferRepo) Store(tr *transfer.Transfer) error {
	log := pgtr.log.WithField("op", "Store").WithField("transfer_id", tr.ID)

	const cmd = `
		INSERT INTO transfers(
			id,
			account_origin_id,
			account_destination_id,
			amount
		) VALUES (
			$1, $2, $3, $4
		)
		RETURNING 
			created_at;
	`

	stmt, _, err := pgtr.stmt(cmd)
	if err != nil {
		log.WithError(err).Error("error on prepare insert statement")
		return err
	}

	row := stmt.QueryRowContext(getContext(),
		tr.ID,
		tr.AccountOriginID,
		tr.AccountDestinationID,
		tr.Amount,
	)

	err = row.Scan(&tr.CreatedAt)
	if err != nil {
		log.WithError(err).Error("error on scan value")
		return err
	}

	err = stmt.Close()
	if err != nil {
		log.WithError(err).Error("error on close statement")
		return err
	}

	return nil
}

func (pgtr *pgTransferRepo) GetByID(transferID string) (tr transfer.Transfer, err error) {
	log := pgtr.log.WithField("op", "GetByID").WithField("transfer_id", transferID)

	const cmd = `
		SELECT
			id,
			account_origin_id,
			account_destination_id,
			amount,
			created_at
		FROM
			transfers
		WHERE
			id = $1`

	stmt, _, err := pgtr.stmt(cmd)
	if err != nil {
		log.WithError(err).Error("error on prepare select statement")
		return
	}

	row := stmt.QueryRowContext(getContext(), transferID)

	err = row.Scan(
		&tr.ID,
		&tr.AccountOriginID,
		&tr.AccountDestinationID,
		&tr.Amount,
		&tr.CreatedAt,
	)
	if err != nil {
		log.WithError(err).Error("error on scan transfer values")
		if err == sql.ErrNoRows {
			err = errors.ErrTransferNotFound
			return
		}
		return

	}

	err = stmt.Close()
	if err != nil {
		log.WithError(err).Error("error on close statement")
		return
	}
	return
}

func (pgtr *pgTransferRepo) ListByAccountID(accountID string) (transfers []transfer.Transfer, err error) {
	log := pgtr.log.WithField("op", "ListByAccountID").WithField("account_id", accountID)

	const cmd = `
		SELECT
			id,
			account_origin_id,
			account_destination_id,
			amount,
			created_at
		FROM
			transfers
		WHERE
			account_origin_id = $1`

	stmt, _, err := pgtr.stmt(cmd)
	if err != nil {
		log.WithError(err).Error("error on prepare select statement")
		return
	}

	rows, err := stmt.QueryContext(getContext(), accountID)
	if err != nil {
		log.WithError(err).Error("error on list transfers")
		if err == sql.ErrNoRows {
			err = errors.ErrTransferNotFound
		}
		return
	}

	transfers = []transfer.Transfer{}

	for rows.Next() {
		var tr transfer.Transfer
		err = rows.Scan(
			&tr.ID,
			&tr.AccountOriginID,
			&tr.AccountDestinationID,
			&tr.Amount,
			&tr.CreatedAt,
		)
		if err != nil {
			log.WithError(err).Error("erros on scan transfer values")
			if err == sql.ErrNoRows {
				err = errors.ErrTransferNotFound
			}
			return
		}

		transfers = append(transfers, tr)
	}

	err = stmt.Close()
	if err != nil {
		log.WithError(err).Error("error on close statement")
		return
	}
	return

}

func (pgtr *pgTransferRepo) GenerateIdentifier() string {
	const cmd = "SELECT gen_random_uuid()"
	var uuid string

	row := pgtr.db.QueryRowContext(getContext(), cmd)
	err := row.Scan(&uuid)
	if err != nil {
		uuid = common.GenUUID()
	}
	return uuid
}

func (pgtr *pgTransferRepo) WithTx(tx driver.Tx) transfer.Repository {
	return &pgTransferRepo{
		db:  pgtr.db,
		log: pgtr.log,
		tx:  tx,
	}
}

func (pgtr *pgTransferRepo) stmt(cmd string) (stmt *sql.Stmt, commit bool, err error) {
	log := pgtr.log.WithField("op", "stmt")

	if pgtr.tx != nil {
		commit = true
		if tx, ok := pgtr.tx.(*sql.Tx); ok {
			stmt, err = tx.PrepareContext(getContext(), cmd)
		} else {
			log.Fatal("tx must be of type *sql.Tx")
		}
		return
	}
	stmt, err = pgtr.db.PrepareContext(getContext(), cmd)
	return
}

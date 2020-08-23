package psql

import (
	"database/sql"
	"database/sql/driver"

	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type pgAccountRepo struct {
	db  *sql.DB
	tx  driver.Tx
	log *logrus.Entry
}

// NewAccountRepository returns a account repository implementation for postgreSQL.
func NewAccountRepository(db *sql.DB) account.Repository {
	return &pgAccountRepo{
		log: common.Logger().WithField("source", "pgAccountRepo"),
		db:  db,
	}
}

func (pgacr *pgAccountRepo) Store(ac *account.Account) error {
	log := pgacr.log.WithField("op", "Store")

	const cmd = `
		INSERT INTO accounts(
			id,
			name,
			cpf,
			secret,
			balance
		) VALUES ($1, $2, $3, $4, $5) RETURNING created_at`

	stmt, _, err := pgacr.stmt(cmd)
	if err != nil {
		log.WithError(err).Error("error on create sql statement")
		return err
	}

	row := stmt.QueryRowContext(
		getContext(),
		ac.ID,
		ac.Name,
		ac.CPF,
		ac.Secret,
		ac.Balance,
	)
	err = row.Scan(&ac.CreatedAt)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23505" {
				return errors.ErrAccountCPFAlreadyExists
			}
		}
		log.WithError(err).Error("error on scan created_at time value")
		return err
	}

	err = stmt.Close()
	if err != nil {
		log.WithError(err).Error("error on close statement")
		return err
	}

	return nil
}

func (pgacr *pgAccountRepo) UpdateBalance(ac account.Account) error {
	log := pgacr.log.WithField("op", "UpdateBalance").WithField("account_id", ac.ID)
	const cmd = `
		UPDATE accounts
			SET balance = $1
		WHERE
			id = $2`

	stmt, _, err := pgacr.stmt(cmd)

	result, err := stmt.Exec(ac.Balance, ac.ID)
	if err != nil {
		log.WithError(err).Error("error on update balance")
		return err
	}

	affecteds, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("error on check affecteds rows")
		return err
	}

	if affecteds == 0 {
		err = errors.ErrAccountNotFound
		log.WithError(err).Error("account not found")
		return err
	}

	err = stmt.Close()
	if err != nil {
		log.WithError(err).Error("error on close statement")
		return err
	}

	return nil
}

func (pgacr *pgAccountRepo) GetByID(accountID string) (ac account.Account, err error) {
	log := pgacr.log.WithField("op", "GetByID").WithField("account_id", accountID)

	const cmd = `
		SELECT
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		FROM 
			accounts
		WHERE
			id = $1
	`

	stmt, _, err := pgacr.stmt(cmd)
	if err != nil {
		log.WithError(err).Error("error on create select account command")
		return
	}

	row := stmt.QueryRow(accountID)
	err = row.Scan(
		&ac.ID,
		&ac.Name,
		&ac.CPF,
		&ac.Secret,
		&ac.Balance,
		&ac.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.ErrAccountNotFound
			log.WithError(err).Error("account not found")
			return
		}
		log.WithError(err).Error("error on scan query values")
		return
	}

	err = stmt.Close()
	if err != nil {
		log.WithError(err).Error("error on close statement")
		return
	}

	return ac, err
}

func (pgacr *pgAccountRepo) GetByCPF(cpf string) (ac account.Account, err error) {
	log := pgacr.log.WithField("op", "GetByCPF").WithField("account_cpf", cpf)

	const cmd = `
		SELECT
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		FROM 
			accounts
		WHERE
			cpf = $1
	`

	stmt, _, err := pgacr.stmt(cmd)
	if err != nil {
		log.WithError(err).Error("error on create select account command")
		return
	}

	row := stmt.QueryRowContext(getContext(), cpf)
	err = row.Scan(
		&ac.ID,
		&ac.Name,
		&ac.CPF,
		&ac.Secret,
		&ac.Balance,
		&ac.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.ErrAccountNotFound
			log.WithError(err).Error("account not found")
			return
		}
		log.WithError(err).Error("error on scan query values")
		return
	}

	err = stmt.Close()
	if err != nil {
		log.WithError(err).Error("error on close statement")
		return
	}

	return
}

func (pgacr *pgAccountRepo) ListAll() (accounts []account.Account, err error) {
	log := pgacr.log.WithField("op", "ListAll")

	const cmd = `
		SELECT
			id,
			name,
			cpf,
			balance,
			created_at
		FROM
			accounts
	`

	stmt, _, err := pgacr.stmt(cmd)
	if err != nil {
		log.WithError(err).Error("error on create statement")
		return
	}

	rows, err := stmt.QueryContext(getContext())
	if err != nil {
		log.WithError(err).Error("error on exec select statement")
		return
	}

	accounts = []account.Account{}

	for rows.Next() {
		var ac account.Account
		err = rows.Scan(
			&ac.ID,
			&ac.Name,
			&ac.CPF,
			&ac.Balance,
			&ac.CreatedAt,
		)
		if err != nil {
			log.WithError(err).Error("erro on scan row values")
			if err == sql.ErrNoRows {
				err = errors.ErrAccountNotFound
			}
			return
		}
		accounts = append(accounts, ac)
	}

	err = stmt.Close()
	if err != nil {
		log.WithError(err).Error("error on close statement")
		return
	}

	return accounts, err
}

func (pgacr *pgAccountRepo) GenerateIdentifier() string {
	const cmd = "SELECT gen_random_uuid()"
	var uuid string

	row := pgacr.db.QueryRowContext(getContext(), cmd)
	err := row.Scan(&uuid)
	if err != nil {
		pgacr.log.Warning("postgresql gen_random_uuid failed, using golang alternative")
		uuid = common.GenUUID()
	}
	return uuid
}

func (pgacr *pgAccountRepo) WithTx(tx driver.Tx) account.Repository {
	return &pgAccountRepo{
		db:  pgacr.db,
		log: pgacr.log,
		tx:  tx,
	}
}

func (pgacr *pgAccountRepo) stmt(cmd string) (stmt *sql.Stmt, commit bool, err error) {
	log := pgacr.log.WithField("op", "stmt")

	if pgacr.tx != nil {
		commit = true
		if tx, ok := pgacr.tx.(*sql.Tx); ok {
			stmt, err = tx.PrepareContext(getContext(), cmd)
		} else {
			log.Fatal("tx must be of type *sql.Tx")
		}
		return
	}

	stmt, err = pgacr.db.PrepareContext(getContext(), cmd)
	return
}

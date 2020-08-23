package pgql

import (
	"database/sql"
	"database/sql/driver"

	_ "github.com/jackc/pgx/v4/stdlib" // register postgresql driver for sql standard lib.
	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type pgRegistry struct {
	db  *sql.DB
	log *logrus.Entry
}

// Connect to database.
func Connect(conStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", conStr)
	if err != nil {
		return db, err
	}

	err = db.PingContext(getContext())
	if err != nil {
		return db, err
	}

	return db, nil
}

// NewRepositoryRegistry returns a new domain.RepositoryRegistry implementation for PostgreSQL.
func NewRepositoryRegistry(db *sql.DB) domain.RepositoryRegistry {
	return &pgRegistry{
		log: common.Logger().WithField("source", "pgRegistry"),
		db:  db,
	}
}

func (pgr *pgRegistry) AccountRepository() account.Repository {
	return NewAccountRepository(pgr.db)
}

func (pgr *pgRegistry) TransferRepository() transfer.Repository {
	return NewTransferRepository(pgr.db)
}

func (pgr *pgRegistry) Tx() driver.Tx {
	tx, err := pgr.db.Begin()
	if err != nil {
		pgr.log.Fatal("fatal error on create transaction")
	}
	return tx
}

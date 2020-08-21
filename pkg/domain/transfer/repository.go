package transfer

import "database/sql/driver"

// Repository manages transfer persistence.
type Repository interface {
	Store(*Transfer) error
	GetByID(string) (Transfer, error)
	ListByAccountID(string) ([]Transfer, error)
	GenerateIdentifier() string
	WithTx(driver.Tx) Repository
}

package account

import "database/sql/driver"

// Repository manages account persistence.
type Repository interface {
	Store(*Account) error
	UpdateBalance(Account) error
	GetByID(string) (Account, error)
	GetByCPF(string) (Account, error)
	ListAll() ([]Account, error)
	GenerateIdentifier() string
	WithTx(driver.Tx) Repository
}

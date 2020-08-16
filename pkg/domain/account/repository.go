package account

// Repository manages account persistence.
type Repository interface {
	Store(Account) error
	UpdateBalance(Account) error
	GetByID(string) (Account, error)
	ListAll() ([]Account, error)
	GenerateIdentifier() string
}

package transfer

// Repository manages transfer persistence.
type Repository interface {
	Store(*Transfer) error
	GetById(string) (Transfer, error)
	ListByAccountID(string) ([]Transfer, error)
	GenerateIndetifier() string
}

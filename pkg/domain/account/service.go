package account

// Service encapsulates domain operations over the account model as a service.
type Service interface {
	CreateAccount(*Account) error
	UpdateBalance(*Account) error
	AccountBalance(string) (Account, error)
	ListAccounts() ([]Account, error)
}

// ValidationService is responsible for validating account data.
type ValidationService interface {
	ValidateForCreation(*Account) error
	ValidateForTransferOrigin(*Account) error
}

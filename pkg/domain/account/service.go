package account

// Service encapsulates domain operations over the account model as a service.
type Service interface {
	CreateAccount(*Account) error
	UpdateBalance(*Account) error
	AccountBalance(string) (BalanceValue, error)
	GetAccount(string) (Account, error)
	ListAccounts() ([]Account, error)
}

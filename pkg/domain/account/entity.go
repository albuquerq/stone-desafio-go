package account

import "time"

//Account is a bank account.
type Account struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
	// Balance in Brazilian real cents BLR.
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate using a strategy.
func (ac Account) Validate(vs ValidationStrategy) error {
	return vs.Validate(ac)
}

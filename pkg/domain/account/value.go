package account

// InputValue to create a Account.
type InputValue struct {
	Name   string `json:"name"`
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
	// Balance in Brazilian real cents BLR.
	Balance int `json:"balance"`
}

// Validate validates using a strategy.
func (icv InputValue) Validate(vs InputCreationValueValidationStrategy) error {
	return vs.Validate(icv)
}

// BalanceValue to get the account balance.
type BalanceValue struct {
	Balance int `json:"balance"`
}

package access

// Credential is value object for access credentials.
type Credential struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

// Description is the value object that describes account ownership information.
type Description struct {
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	AccountID string `json:"account_id"`
}

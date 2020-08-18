package access

// Credential is value object for access credentials.
type Credential struct {
	CPF    string
	Secret string
}

// Description is the value object that describes account ownership information.
type Description struct {
	Name      string
	CPF       string
	AccountID string
}

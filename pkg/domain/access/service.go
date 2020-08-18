package access

// Service responsible for authentication.
type Service interface {
	Authenticate(Credential) (Description, error)
}

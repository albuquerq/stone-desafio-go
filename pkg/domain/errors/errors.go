package errors

import (
	"errors"
	"fmt"
)

// Domain errors.
var (
	ErrDomain                   = errors.New("domain error")
	ErrInsufficientBalance      = fmt.Errorf("%w: %v", ErrDomain, "insufficient balance")
	ErrInvalidAccessCredentials = fmt.Errorf("%w: %v", ErrDomain, "invalid access credentials")
	ErrNoHasUniqueIdentity      = fmt.Errorf("%w: %v", ErrDomain, "no has unique identity")
	ErrAccountNotFound          = fmt.Errorf("%w: %v", ErrDomain, "account not found")
)

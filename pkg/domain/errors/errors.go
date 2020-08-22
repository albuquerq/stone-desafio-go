package errors

import (
	"errors"
	"fmt"
)

// Domain errors.
var (
	ErrDomain                             = errors.New("domain error")
	ErrInvalidAccessCredentials           = fmt.Errorf("%w: %v", ErrDomain, "invalid access credentials")
	ErrNoHasUniqueIdentity                = fmt.Errorf("%w: %v", ErrDomain, "no has unique identity")
	ErrAccountNotFound                    = fmt.Errorf("%w: %v", ErrDomain, "account not found")
	ErrTransferNotFound                   = fmt.Errorf("%w: %v", ErrDomain, "transfer not found")
	ErrTransferNotAllowed                 = fmt.Errorf("%w: %v", ErrDomain, "transfer not allowed")
	ErrTransferInsufficientBalance        = fmt.Errorf("%w: %v", ErrTransferNotAllowed, "the origin account has insufficient balance")
	ErrTransferMissingData                = fmt.Errorf("%w: %v", ErrTransferNotAllowed, "missing data")
	ErrTransferMissingAmount              = fmt.Errorf("%w: %v", ErrTransferMissingData, "missing amount")
	ErrTransferBetweenSameAccount         = fmt.Errorf("%w: %v", ErrTransferNotAllowed, "transfer between the same account")
	ErrTransferAccountOriginNotFound      = fmt.Errorf("%w: %v", ErrTransferNotAllowed, "oringin account not found")
	ErrTransferAccountDestinationNotFound = fmt.Errorf("%w: %v", ErrTransferNotAllowed, "destination account not found")
)

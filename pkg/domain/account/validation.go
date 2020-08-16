package account

// ValidationStrategy for data validation.
type ValidationStrategy interface {
	Validate(Account) error
}

// InputCreationValueValidationStrategy for data validation.
type InputCreationValueValidationStrategy interface {
	Validate(InputValue) error
}

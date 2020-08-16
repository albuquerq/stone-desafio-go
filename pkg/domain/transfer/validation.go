package transfer

// ValidationStrategy for data validation.
type ValidationStrategy interface {
	Validate(Transfer) error
}

// InputValidationStrategy for data validation.
type InputValidationStrategy interface {
	Validate(InputValue) error
}

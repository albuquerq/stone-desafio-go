package transfer

// Service encapsulates domain operations over the transfer model as a service.
type Service interface {
	Transfer(accountOringinID string, accountDestinationID string, amount int) (Transfer, error)
	ListTransfersFromAccount(accountID string) ([]Transfer, error)
}

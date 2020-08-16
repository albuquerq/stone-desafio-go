package transfer

// InputValue to create a transfer.
type InputValue struct {
	AccountOriginID      string `json:"account_origin_id"`
	AccountDestinationID string `json:"account_destination_id"`
	// Amount in Brazilian real cents BLR..
	Amount int `json:"amount"`
}

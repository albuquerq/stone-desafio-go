package transfer

import "time"

// Transfer is a transfer between two accounts.
type Transfer struct {
	ID                   string `json:"id"`
	AccountOriginID      string `json:"account_origin_id"`
	AccountDestinationID string `json:"account_destination_id"`
	// Amount in Brazilian real cents BLR..
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate validates using a validation strategy.
func (t Transfer) Validate(vs ValidationStrategy) error {
	return vs.Validate(t)
}

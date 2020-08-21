package mem

// // Transactions not applicable to in-memory storage.
type memTx struct{}

func (memTx) Commit() error {
	return nil
}

func (memTx) Rollback() error {
	return nil
}

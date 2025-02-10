package base

// FailedTx returns an empty string and the given error.
// It is intended to be used to signal that a transaction has failed.
func FailedTx(err error) (string, error) {
	return "", err
}

// SuccessTx returns the given hash and a nil error.
// It is intended to be used to signal that a transaction has succeeded.
func SuccessTx(hash string) (string, error) {
	return hash, nil
}

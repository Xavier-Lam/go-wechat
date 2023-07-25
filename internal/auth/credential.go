package auth

// It would be much better if Go supports covariance...
type CredentialManager interface {
	// Get the latest credential
	Get() (interface{}, error)

	// Set the latest credential
	Set(credential interface{}) error

	// Renew credential
	Renew() (interface{}, error)

	// Delete a credential
	Delete() error
}

package auth

type CredentialManager[T interface{}] interface {
	// Get the latest credential
	Get() (*T, error)

	// Set the latest credential
	Set(credential *T) error

	// Renew credential
	Renew() (*T, error)

	// Delete a credential
	Delete() error
}

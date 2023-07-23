package auth

import "errors"

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

type authCredentialManager struct {
	auth Auth
}

// Provide `Auth`
func NewAuthCredentialManager(auth Auth) CredentialManager {
	return &authCredentialManager{auth: auth}
}

func (cm *authCredentialManager) Get() (interface{}, error) {
	if cm.auth == nil {
		return errors.New("auth not set"), nil
	}
	return cm.auth, nil
}

func (cm *authCredentialManager) Set(credential interface{}) error {
	return errors.New("not settable")
}

func (cm *authCredentialManager) Renew() (interface{}, error) {
	return nil, errors.New("not renewable")
}

func (cm *authCredentialManager) Delete() error {
	return errors.New("not deletable")
}

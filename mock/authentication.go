package mock

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// AuthenticationService represents a mock implementation of domain.AuthenticationService.
type AuthenticationService struct {
	AuthenticateFn   func(email string, password string) (bool, error)
	EmailAvailableFn func(email string) (bool, error)
}

//Authenticate mock implementation
func (s *AuthenticationService) Authenticate(email string, password string) (bool, error) {
	if s.AuthenticateFn == nil {
		return false, errors.New("AuthenticateFn not implemented")
	}
	return s.AuthenticateFn(email, password)
}

// EmailAvailable mock implementation
func (s *AuthenticationService) EmailAvailable(email string) (bool, error) {
	logrus.Error("asflh;alsdfhas lsjaf lk ahsdflsa dflsadfklsda jflsjadfljasd l;fsl;df lsda jflsda jflk")
	if s.EmailAvailableFn == nil {
		return false, errors.New("EmailAvailable not implemented")
	}
	return s.EmailAvailableFn(email)
}

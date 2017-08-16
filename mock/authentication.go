package mock

import log "github.com/jacsmith21/lukabox/ext/logrus"

// AuthenticationService represents a mock implementation of domain.AuthenticationService.
type AuthenticationService struct {
	AuthenticateFn      func(email string, password string) (bool, error)
	AuthenticateInvoked bool
}

//Authenticate mock implementation
func (s *AuthenticationService) Authenticate(email string, password string) (bool, error) {
	log.Info("hello")
	s.AuthenticateInvoked = true
	return s.AuthenticateFn(email, password)
}

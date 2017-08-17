package mock

// AuthenticationService represents a mock implementation of domain.AuthenticationService.
type AuthenticationService struct {
	AuthenticateFn      func(email string, password string) (bool, error)
	AuthenticateInvoked bool

	EmailAvailableFn      func(email string) (bool, error)
	EmailAvailableInvoked bool
}

//Authenticate mock implementation
func (s *AuthenticationService) Authenticate(email string, password string) (bool, error) {
	s.AuthenticateInvoked = true
	return s.AuthenticateFn(email, password)
}

// EmailAvailable mock implementation
func (s *AuthenticationService) EmailAvailable(email string) (bool, error) {
	s.EmailAvailableInvoked = true
	return s.EmailAvailableFn(email)
}

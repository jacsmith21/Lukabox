package mock

// AuthenticationService represents a mock implementation of domain.AuthenticationService.
type AuthenticationService struct {
	AuthenticateFn      func(email string, password string) (bool, error)
	AuthenticateInvoked bool
}

//Authenticate mock implementation
func (s *AuthenticationService) Authenticate(email string, password string) (bool, error) {
	s.AuthenticateInvoked = true
	return s.AuthenticateFn(email, password)
}

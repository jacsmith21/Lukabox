package db

//AuthenticationService AuthenticationService implementation
type AuthenticationService struct {
}

//Authenticate authenticates a user with credentials
func (s *AuthenticationService) Authenticate(email string, password string) bool {
	for _, u := range users {
		if u.Email == email {
			if u.Password == password {
				return true
			}
			return false
		}
	}
	return false
}

package db

//AuthenticationService AuthenticationService implementation
type AuthenticationService struct {
}

//Authenticate authenticates a user with credentials
func (s *AuthenticationService) Authenticate(email string, password string) (bool, error) {
	for _, u := range users {
		if u.Email == email {
			if u.Password == password {
				return true, nil
			}
			return false, nil
		}
	}
	return false, nil
}

// EmailAvailable checks email availability
func (s *AuthenticationService) EmailAvailable(email string) (bool, error) {
	for _, u := range users {
		if u.Email == email {
			return false, nil
		}
	}
	return true, nil
}

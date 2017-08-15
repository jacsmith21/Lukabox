package domain

//Credentials a reguler user credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//AuthenticationService credentials services
type AuthenticationService interface {
	Authenticate(email string, password string) (bool, error)
}

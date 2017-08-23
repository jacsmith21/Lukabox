package domain

//Credentials a reguler user credentials
type Credentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

//Token jwt token
type Token struct {
	Token string `json:"token"`
}

//AuthenticationService credentials services
type AuthenticationService interface {
	Authenticate(email string, password string) (bool, error)
	EmailAvailable(email string) (bool, error)
}

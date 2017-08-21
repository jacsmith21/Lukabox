package domain

//User a reguler user
type User struct {
	ID        int    `json:"id"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Archived  bool   `json:"archived"`
}

//UserService database services
type UserService interface {
	UserByID(id int) (*User, error)
	UserByEmail(email string) (*User, error)
	Users() ([]*User, error)
	ValidateUser(user *User) error
	InsertUser(user *User) error
	UpdateUser(id int, user *User) error
}

package domain

//User a reguler user
type User struct {
	ID        int    `json:"id"`
	Password  string `json:"password" validate:"required"`
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Archived  bool   `json:"archived"`
}

//UserService database services
type UserService interface {
	UserByID(id int) (*User, error)
	UserByEmail(email string) (*User, error)
	Users() ([]*User, error)
	InsertUser(user *User) error
	UpdateUser(id int, user *User) error
}

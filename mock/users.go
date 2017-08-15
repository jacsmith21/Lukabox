package mock

import "github.com/jacsmith21/lukabox/domain"

// UserService represents a mock implementation of myapp.UserService.
type UserService struct {
	UserByIDFn      func(id int) (*domain.User, error)
	UserByIDInvoked bool

	UserByEmailFn      func(email string) (*domain.User, error)
	UserByEmailInvoked bool

	UsersFn      func() ([]*domain.User, error)
	UsersInvoked bool

	CreateUserFn      func(user *domain.User) error
	CreateUserInvoked bool

	UpdateUserFn      func(id int, user *domain.User) error
	UpdateUserInvoked bool

	AuthenticateUserFn      func(email string, password string) (bool, error)
	AuthenticateUserInvoked bool
}

//UserByID invokes the mock implementation and marks the function as invoked.
func (s *UserService) UserByID(id int) (*domain.User, error) {
	s.UserByIDInvoked = true
	return s.UserByIDFn(id)
}

//UserByEmail mock implementation
func (s *UserService) UserByEmail(email string) (*domain.User, error) {
	s.UserByEmailInvoked = true
	return s.UserByEmailFn(email)
}

//Users mock implementation
func (s *UserService) Users() ([]*domain.User, error) {
	s.UsersInvoked = true
	return s.UsersFn()
}

//CreateUser mock implementation
func (s *UserService) CreateUser(user *domain.User) error {
	s.CreateUserInvoked = true
	return s.CreateUserFn(user)
}

//UpdateUser mock implementation
func (s *UserService) UpdateUser(id int, user *domain.User) error {
	s.UpdateUserInvoked = true
	return s.UpdateUserFn(id, user)
}

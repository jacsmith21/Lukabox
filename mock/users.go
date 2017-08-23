package mock

import (
	"errors"

	"github.com/jacsmith21/lukabox/domain"
)

// UserService represents a mock implementation of domain.UserService.
type UserService struct {
	UserByIDFn    func(id int) (*domain.User, error)
	UserByEmailFn func(email string) (*domain.User, error)
	UsersFn       func() ([]*domain.User, error)
	InsertUserFn  func(user *domain.User) error
	UpdateUserFn  func(id int, user *domain.User) error
}

//UserByID invokes the mock implementation and marks the function as invoked.
func (s *UserService) UserByID(id int) (*domain.User, error) {
	if s.UserByIDFn == nil {
		return nil, errors.New("UserByIDFn not implemented")
	}
	return s.UserByIDFn(id)
}

//UserByEmail mock implementation
func (s *UserService) UserByEmail(email string) (*domain.User, error) {
	if s.UserByEmailFn == nil {
		return nil, errors.New("UserByEmailFn not implemented")
	}
	return s.UserByEmailFn(email)
}

//Users mock implementation
func (s *UserService) Users() ([]*domain.User, error) {
	if s.UsersFn == nil {
		return nil, errors.New("UsersFn not implemented")
	}
	return s.UsersFn()
}

// InsertUser mock implementation
func (s *UserService) InsertUser(user *domain.User) error {
	if s.InsertUserFn == nil {
		return errors.New("InsertUserFn not implemeted")
	}
	return s.InsertUserFn(user)
}

//UpdateUser mock implementation
func (s *UserService) UpdateUser(id int, user *domain.User) error {
	if s.UpdateUserFn == nil {
		return errors.New("UpdateUserFn not implemented")
	}
	return s.UpdateUserFn(id, user)
}

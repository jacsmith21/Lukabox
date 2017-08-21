package mock

import (
	"errors"

	"github.com/jacsmith21/lukabox/domain"
)

// UserService represents a mock implementation of domain.UserService.
type UserService struct{}

//UserByID invokes the mock implementation and marks the function as invoked.
func (s *UserService) UserByID(id int) (*domain.User, error) {
	if id == 1 {
		user := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
		return &user, nil
	} else if id == 2 {
		user := domain.User{ID: 2, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
		return &user, nil
	} else if id == 3 {
		return nil, errors.New("test error")
	}
	return nil, nil
}

//UserByEmail mock implementation
func (s *UserService) UserByEmail(email string) (*domain.User, error) {
	if email == "jacob.smith@unb.ca" {
		user := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
		return &user, nil
	} else if email == "" {
		return nil, errors.New("no email supplied")
	}
	return nil, nil
}

var usersCount = 0

//Users mock implementation
func (s *UserService) Users() ([]*domain.User, error) {
	usersCount++
	if usersCount == 1 {
		users := []*domain.User{
			{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
		}
		return users, nil
	} else if usersCount == 2 {
		return nil, nil
	}
	return nil, errors.New("test error")
}

// ValidateUser mock implementation
func (s *UserService) ValidateUser(user *domain.User) error {
	return nil
}

var insertUserCount = 0

// InsertUser mock implementation
func (s *UserService) InsertUser(user *domain.User) error {
	if insertUserCount == 0 {
		return nil
	}
	insertUserCount++
	return errors.New("test error")
}

//UpdateUser mock implementation
func (s *UserService) UpdateUser(id int, user *domain.User) error {
	if id != 1 {
		return errors.New("expected id to be 1")
	}
	return nil
}

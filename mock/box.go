package mock

import (
	"errors"

	"github.com/jacsmith21/lukabox/domain"
)

// BoxService mock implementation
type BoxService struct {
	BoxesFn            func() ([]*domain.Box, error)
	BoxFn              func(userID int, id int) (*domain.Box, error)
	InsertBoxFn        func(box *domain.Box) error
	UpdateBoxFn        func(id int, box *domain.Box) error
	InsertOpenEventFn  func(openEvent *domain.OpenEvent) error
	InsertCloseEventFn func(closeEvent *domain.CloseEvent) error
}

// Boxes mock implementation
func (s *BoxService) Boxes(userID int) ([]*domain.Box, error) {
	if s.BoxesFn == nil {
		return nil, errors.New("Boxes not implemented")
	}
	return s.BoxesFn()
}

// Box mock implementation
func (s *BoxService) Box(userID int, id int) (*domain.Box, error) {
	if s.BoxFn == nil {
		return nil, errors.New("BoxFn not implemented")
	}
	return s.BoxFn(userID, id)
}

// InsertBox mock implementation
func (s *BoxService) InsertBox(box *domain.Box) error {
	if s.InsertBoxFn == nil {
		return errors.New("InsertBoxFn not implemented")
	}
	return s.InsertBoxFn(box)
}

// UpdateBox mock implementation
func (s *BoxService) UpdateBox(id int, box *domain.Box) error {
	if s.UpdateBoxFn == nil {
		return errors.New("UpdateBoxFn not implemented")
	}
	return s.UpdateBoxFn(id, box)
}

// InsertOpenEvent mock implementation
func (s *BoxService) InsertOpenEvent(openEvent *domain.OpenEvent) error {
	if s.InsertOpenEventFn == nil {
		return errors.New("InsertOpenEventFn not implemented")
	}
	return s.InsertOpenEventFn(openEvent)
}

// InsertCloseEvent mock implementation
func (s *BoxService) InsertCloseEvent(closeEvent *domain.CloseEvent) error {
	if s.InsertCloseEventFn == nil {
		return errors.New("InsertCloseEventFn not implemented")
	}
	return s.InsertCloseEventFn(closeEvent)
}

package mock

import (
	"errors"

	"github.com/jacsmith21/lukabox/domain"
)

// BoxService mock implementation
type BoxService struct {
	BoxesFn     func() ([]*domain.Box, error)
	BoxByIDFn   func(id int) (*domain.Box, error)
	InsertBoxFn func(box *domain.Box) error
	UpdateBoxFn func(id int, box *domain.Box) error
}

// Boxes mock implementation
func (s *BoxService) Boxes() ([]*domain.Box, error) {
	if s.BoxesFn == nil {
		return nil, errors.New("Boxes not implemented")
	}
	return s.BoxesFn()
}

// BoxByID mock implementation
func (s *BoxService) BoxByID(id int) (*domain.Box, error) {
	if s.BoxByIDFn == nil {
		return nil, errors.New("BoxByIDFn not implemented")
	}
	return s.BoxByIDFn(id)
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

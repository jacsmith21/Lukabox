package mock

import (
	"errors"

	"github.com/jacsmith21/lukabox/domain"
)

// PillService represents a mock implementation of domain.PillService.
type PillService struct {
	PillFn       func(id int) (*domain.Pill, error)
	PillsFn      func(id int) ([]*domain.Pill, error)
	CreatePillFn func(pill *domain.Pill) error
	UpdatePillFn func(id int, pill *domain.Pill) error
}

//Pill invokes the mock implementation and marks the function as invoked.
func (s *PillService) Pill(id int) (*domain.Pill, error) {
	if s.PillFn == nil {
		return nil, errors.New("PillFn not implemented")
	}
	return s.PillFn(id)
}

//Pills mock implementation
func (s *PillService) Pills(id int) ([]*domain.Pill, error) {
	if s.PillsFn == nil {
		return nil, errors.New("PillsFn not implemented")
	}
	return s.PillsFn(id)
}

//CreatePill mock implementation
func (s *PillService) CreatePill(pill *domain.Pill) error {
	if s.CreatePillFn == nil {
		return errors.New("CreatePillFn not implemented")
	}
	return s.CreatePillFn(pill)
}

//UpdatePill mock implementation
func (s *PillService) UpdatePill(id int, pill *domain.Pill) error {
	if s.UpdatePillFn == nil {
		return errors.New("UpdatePillFn not implemented")
	}
	return s.UpdatePillFn(id, pill)
}

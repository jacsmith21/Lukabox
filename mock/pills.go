package mock

import "github.com/jacsmith21/lukabox/domain"

// PillService represents a mock implementation of domain.PillService.
type PillService struct {
	PillFn      func(id int) (*domain.Pill, error)
	PillInvoked bool

	PillsFn      func(id int) ([]*domain.Pill, error)
	PillsInvoked bool

	CreatePillFn      func(pill *domain.Pill) error
	CreatePillInvoked bool

	UpdatePillFn      func(id int, pill *domain.Pill) error
	UpdatePillInvoked bool
}

//Pill invokes the mock implementation and marks the function as invoked.
func (s *PillService) Pill(id int) (*domain.Pill, error) {
	s.PillInvoked = true
	return s.PillFn(id)
}

//Pills mock implementation
func (s *PillService) Pills(id int) ([]*domain.Pill, error) {
	s.PillsInvoked = true
	return s.PillsFn(id)
}

//CreatePill mock implementation
func (s *PillService) CreatePill(pill *domain.Pill) error {
	s.CreatePillInvoked = true
	return s.CreatePillFn(pill)
}

//UpdatePill mock implementation
func (s *PillService) UpdatePill(id int, pill *domain.Pill) error {
	s.UpdatePillInvoked = true
	return s.UpdatePillFn(id, pill)
}

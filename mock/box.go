package mock

import (
	"errors"

	"github.com/jacsmith21/lukabox/domain"
)

// BoxService mock implementation
type BoxService struct {
	InsertOpenEventFn  func(openEvent *domain.OpenEvent) error
	InsertCloseEventFn func(closeEvent *domain.CloseEvent) error
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

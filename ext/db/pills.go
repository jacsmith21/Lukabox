package db

import (
	"errors"
	"time"

	"github.com/jacsmith21/lukabox/domain"
)

// PillService implementation of domain.PillService
type PillService struct {
}

var pills = []*domain.Pill{
	{ID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{time.Now()}, Archived: false},
}

//CreatePill creates a pill in the database
func (s *PillService) CreatePill(pill *domain.Pill) error {
	pill.ID = pills[len(pills)-1].ID + 1
	pills = append(pills, pill)
	return nil
}

//Pill retrieves a pill from the database
func (s *PillService) Pill(id int) (*domain.Pill, error) {
	for _, p := range pills {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("pill not found")
}

//Pills retrieves a user's pills from the database
func (s *PillService) Pills(id int) ([]*domain.Pill, error) {
	userPills := []*domain.Pill{}
	for _, p := range pills {
		if p.UserID == id {
			userPills = append(userPills, p)
		}
	}
	return userPills, nil
}

//UpdatePill updates a pill in the datbase
func (s *PillService) UpdatePill(id int, pill *domain.Pill) error {
	for i, p := range pills {
		if p.ID == id {
			pills[i] = pill
			return nil
		}
	}
	return errors.New("pill not found")
}

package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/jacsmith21/lukabox/domain"
)

var pills = []*domain.Pill{
	{ID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{time.Now()}, Archived: false},
}

//CreatePill creates a pill in the database
func CreatePill(pill *domain.Pill) (string, error) {
	pill.ID = pills[len(pills)-1].ID + 1
	pills = append(pills, pill)
	return fmt.Sprintf("%d", pill.ID), nil
}

//GetPills retrieves a pill from the database
func GetPills() ([]*domain.Pill, error) {
	return pills, nil
}

//GetPill retrieves a pill from the database
func GetPill(id int) (*domain.Pill, error) {
	for _, p := range pills {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("pill not found")
}

//GetPillsByUser retrieves a user's pills from the database
func GetPillsByUser(id int) ([]*domain.Pill, error) {
	userPills := []*domain.Pill{}
	for _, p := range pills {
		if p.UserID == id {
			userPills = append(userPills, p)
		}
	}
	return userPills, nil
}

//UpdatePill updates a pill in the datbase
func UpdatePill(id int, pill *domain.Pill) (*domain.Pill, error) {
	for i, p := range pills {
		if p.ID == id {
			pills[i] = pill
			return p, nil
		}
	}
	return nil, errors.New("pill not found")
}

package domain

import "time"

//Pill a pill or other form of medication
type Pill struct {
	PillID     int         `json:"pillId"`
	UserID     int         `json:"id"`
	Name       string      `json:"name"`
	DaysOfWeek []int       `json:"daysOfWeek"`
	TimesOfDay []time.Time `json:"timesOfDay"`
	Archived   bool        `json:"archived"`
}

//PillService database services
type PillService interface {
	Pill(id int) (*Pill, error)
	Pills(id int) ([]*Pill, error)
	CreatePill(pill *Pill) error
	UpdatePill(id int, pill *Pill) error
}

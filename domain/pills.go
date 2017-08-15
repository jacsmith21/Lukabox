package domain

import "time"

//Pill a pill or other form of medication
type Pill struct {
	UserID     int         `json:"userId"`
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	DaysOfWeek []int       `json:"daysOfWeek"`
	TimesOfDay []time.Time `json:"timesOfDay"`
	Archived   bool        `json:"archived"`
}

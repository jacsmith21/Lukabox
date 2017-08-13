package domain

import "time"

//Pill a pill or other form of medication
type Pill struct {
	User          User        `json:"user"`
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	DaysOfWeek    []int       `json:"daysOfWeek"`
	TimesOfDay    []time.Time `json:"timesOfDay"`
	DateTimeAdded time.Time   `json:"dateTimeJoined"`
}

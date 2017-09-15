package domain

import "time"

// Box a box
type Box struct {
	ID     int
	UserID int
}

// OpenEvent an opening event
type OpenEvent struct {
	ID     int
	CompID int
	UserID int
	Time   time.Time
}

// CloseEvent a closing event
type CloseEvent struct {
	ID     int       `json:"id"`
	CompID int       `json:"compId"`
	UserID int       `json:"userId" validate:"required"`
	Time   time.Time `json:"time" validate:"required"`
}

// BoxService database service
type BoxService interface {
	InsertOpenEvent(openEvent *OpenEvent) error
	InsertCloseEvent(closeEvent *CloseEvent) error
}

package domain

import "time"

// Box a box
type Box struct {
	ID     int
	UserID int
}

// OpenEvent an opening event
type OpenEvent struct {
	ID     int       `json:"id"`
	CompID int       `json:"compId" validate:"required"`
	UserID int       `json:"userId"`
	Time   time.Time `json:"time" validate:"required"`
}

// CloseEvent a closing event
type CloseEvent struct {
	ID     int       `json:"id"`
	CompID int       `json:"compId" validate:"required"`
	UserID int       `json:"userId"`
	Time   time.Time `json:"time" validate:"required"`
}

// BoxService database service
type BoxService interface {
	InsertOpenEvent(openEvent *OpenEvent) error
	InsertCloseEvent(closeEvent *CloseEvent) error
}

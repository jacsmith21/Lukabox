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
	BoxID  int
	Time   time.Time
}

// CloseEvent a closing event
type CloseEvent struct {
	ID     int
	CompID int
	BoxID  int
	Time   time.Time
}

// BoxService database service
type BoxService interface {
	Boxes(userID int) ([]*Box, error)
	Box(userID int, ID int) (*Box, error)
	InsertBox(box *Box) error
	UpdateBox(id int, box *Box) error
	InsertOpenEvent(opendEvent *OpenEvent) error
	InsertCloseEvent(opendEvent *CloseEvent) error
}

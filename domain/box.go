package domain

import "time"

// Box a box
type Box struct {
	BoxID  int
	UserID int
}

// Open an opening event
type Open struct {
	OpenID int
	BoxID  int
	Time   time.Time
}

// BoxService database service
type BoxService interface {
	Boxes() ([]*Box, error)
	BoxByID(id int) (*Box, error)
	InsertBox(box *Box) error
	UpdateBox(id int, box *Box) error
}

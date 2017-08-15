package api

import (
	"testing"

	"github.com/jacsmith21/lukabox/mock"
)

func TestPills(t *testing.T) {
	var ps mock.PillService
	var pa PillAPI
	pa.PillService = &ps
}

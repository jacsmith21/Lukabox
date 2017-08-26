package stc

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
)

// OpenEventRequest request structure
type OpenEventRequest struct {
	*domain.OpenEvent
}

// Bind post-processing
func (s *OpenEventRequest) Bind(r *http.ResponseWriter) error {
	return nil
}

// OpenEventResponse reponse structure
type OpenEventResponse struct {
	*domain.OpenEvent
}

// Render pre-processing
func (s *OpenEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewOpenEventResponse open event response
func NewOpenEventResponse(openEvent *domain.OpenEvent) render.Renderer {
	return &OpenEventResponse{OpenEvent: openEvent}
}

// NewOpenEventListReponse open event list response
func NewOpenEventListReponse(openEvenets []*domain.OpenEvent) []render.Renderer {
	list := []render.Renderer{}
	for _, openEvent := range openEvenets {
		list = append(list, NewOpenEventResponse(openEvent))
	}
	return list
}

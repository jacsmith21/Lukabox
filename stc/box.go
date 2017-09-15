package stc

import (
	"errors"
	"net/http"

	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/ext/render"
)

// OpenEventRequest request structure
type OpenEventRequest struct {
	*domain.OpenEvent
}

// Bind post-processing
func (e *OpenEventRequest) Bind(r *http.Request) error {
	tmp := r.Context().Value("user")
	if tmp == nil {
		return errors.New("no user in open event request context")
	}
	user := tmp.(*domain.User)
	e.UserID = user.ID
	return nil
}

// CloseEventRequest CloseEventRequest
type CloseEventRequest struct {
	*domain.CloseEvent
}

// Bind post-processing
func (e *CloseEventRequest) Bind(r *http.Request) error {
	tmp := r.Context().Value("user")
	if tmp == nil {
		return errors.New("no user in close event request context")
	}
	user := tmp.(*domain.User)
	e.UserID = user.ID
	return nil
}

// OpenEventResponse reponse structure
type OpenEventResponse struct {
	*domain.OpenEvent
}

// Render pre-processing
func (e *OpenEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
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

package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
)

//PillsByUser returns the pills associated with the user
func PillsByUser(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	w.Write([]byte(fmt.Sprintf("welcome :) it worked! id: %v", claims["id"])))
}

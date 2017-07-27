package api

import (
  "net/http"
  "fmt"
  "errors"

  "github.com/go-chi/render"
)

//User a reguler user
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Archived  bool   `json:"archived"`
}

var users = []*User{
	{ID: 1, Email: "jacob.smith@unb.ca", FirstName: "Jacob", LastName: "Smith", Archived: false},
	{ID: 2, Email: "j.a.smith@live.ca", FirstName: "Jacob", LastName: "Smith", Archived: false},
	{ID: 3, Email: "jacobsmithunb@gmail.com", FirstName: "Jacob", LastName: "Smith", Archived: false},
}

// UserResponse for json
type UserResponse struct {
	*User
}

// Render does pre-processing before a response is marshalled
func (rd *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewUserResponse ceates a new user reponse
func NewUserResponse(user *User) *UserResponse {
	resp := &UserResponse{User: user}
	return resp
}

// NewUserListResponse creates a new renderer list of reponses
func NewUserListResponse(users []*User) []render.Renderer {
	list := []render.Renderer{}
	for _, user := range users {
		list = append(list, NewUserResponse(user))
	}
	return list
}

// ListUsers lists the users using the RenderList function
func ListUsers(w http.ResponseWriter, r *http.Request) {
	err := render.RenderList(w, r, NewUserListResponse(users))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

//UserRequest a reuqest to create a user
type UserRequest struct {
	*User
}

// Bind post-processing after decode
func (u *UserRequest) Bind(r *http.Request) error {
	u.User.Archived = false
	return nil
}

//CreateUser creates a user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	err := render.Bind(r, data)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	user := data.User
	dbCreateUser(user)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserResponse(user))
}

// UpdateUser updates the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	data := &UserRequest{User: user}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user = data.User
	dbUpdateUser(user.ID, user)

	render.Render(w, r, NewUserResponse(user))
}


func dbCreateUser(user *User) (string, error) {
	user.ID = users[len(users)-1].ID + 1
	users = append(users, user)
	return fmt.Sprintf("%d", user.ID), nil
}

func dbGetUser(id int64) (*User, error) {
	for _, u := range users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func dbUpdateUser(id int64, user *User) (*User, error) {
	for i, u := range users {
		if u.ID == id {
			users[i] = user
			return u, nil
		}
	}
	return nil, errors.New("article not found")
}

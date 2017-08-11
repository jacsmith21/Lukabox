package lukabox

//User a reguler user
type User struct {
	ID        int    `json:"id"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Archived  bool   `json:"archived"`
}

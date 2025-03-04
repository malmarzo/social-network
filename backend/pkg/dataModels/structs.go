package datamodels

//Response struct to be sent back
type Response struct {
	Code     int    `json:"code"`
	Status   string `json:"status"`
	ErrorMsg string `json:"error_msg"`
	Group 	 Group	`json:"group"`
}

//User struct to store user data
type User struct {
	ID string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Password string `json:"password"`
	DOB string `json:"dob"`
	Nickname string `json:"nickname"`
	About string `json:"about_me"`
	Avatar string `json:"avatar"`
	Private bool `json:"private"`
	CreatedAt string `json:"created_at"`
}

// group struct to store groups
type Group struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatorID   string       `json:"creator_id"`
    CreatedAt   string `json:"created_at"`
	FirstName       string    `json:"firstname"`
	LastName       string    `json:"lastname"`
}

// Invite to invite a person to join a group
type Invite struct {
	GroupID  int `json:"group_id"`
	UserID   string `json:"user_id"`
	InvitedBy string `json:"invited_by"`
}

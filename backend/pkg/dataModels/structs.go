package datamodels

//Response struct to be sent back
type Response struct {
	Code     int    `json:"code"`
	Status   string `json:"status"`
	ErrorMsg string `json:"error_msg"`
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
package datamodels

// Response struct to be sent back
type Response struct {
	Code     int         `json:"code"`
	Status   string      `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	ErrorMsg string      `json:"error_msg,omitempty"`
}

// User struct to store user data
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	DOB       string `json:"dob"`
	Nickname  string `json:"nickname"`
	About     string `json:"about_me"`
	Avatar    string `json:"avatar"`
	Private   bool   `json:"private"`
	CreatedAt string `json:"created_at"`
}

// ProfileCard struct that returns the number of posts, followers, and following, username, avatar, and extention of the avatar
type ProfileCard struct {
	Nickanme       string `json:"nickname"`
	NumOfPosts     int    `json:"num_of_posts"`
	NumOfFollowers int    `json:"num_of_followers"`
	NumOfFollowing int    `json:"num_of_following"`
	Avatar         string `json:"avatar"`
	AvatarMimeType string `json:"avatar_mime_type"`
}

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

type UserLogin struct {
	UserID       string `json:"user_id"`
	UserNickname string `json:"user_nickname"`
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

// Post struct
type Post struct {
	PostID        string `json:"post_id"`
	UserID        string `json:"user_id"`
	UserNickname  string `json:"user_nickname"`
	PostTitle     string `json:"post_title"`
	Content       string `json:"content"`
	PostPrivacy   string `json:"post_privacy"`
	PostImage     string `json:"post_image"` // Will contain base64 string after conversion
	NumOfLikes    int    `json:"num_of_likes"`
	NumOfDislikes int    `json:"num_of_dislikes"`
	NumOfComments int    `json:"num_of_comments"`
	CreatedAt     string `json:"created_at"`
	AllowedUsers  string `json:"allowed_users"`
	ImageMimeType string `json:"image_mime_type"` // For content-type header
	ImageDataURL  []byte `json:"-"`               // Temporary storage, won't be sent in JSON
}

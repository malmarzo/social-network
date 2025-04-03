package datamodels

//Response struct to be sent back
type Response struct {
	Code     int    `json:"code"`
	Status   string `json:"status"`
	ErrorMsg string `json:"error_msg"`
	Group 	 Group	`json:"group"`
	Users	 []User `json:"users"`
	Data     interface{} `json:"data,omitempty"`
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
	// GroupMembers  []User      `json:"group_members"` 
	 CurrentUser   string       `json:"current_user"` 
	 ChatHistory	[]GroupMessage `json:"chat_history"`
	 EventHistory    []EventMessage   `json:"event_history"`
	 EventResponsesHistory []EventResponseMessage `json:"event_responses_history"`
}

// Invite to invite a person to join a group
type Invite struct {
	GroupID  int `json:"group_id"`
	UserID   string `json:"user_id"`
	InvitedBy string `json:"invited_by"`
	Accepted bool `json:"accepted"`
}


type Request struct {
	GroupID  int `json:"group_id"`
	GroupCreator string `json:"group_creator"`
	UserID 		 string  `json:"user_id"`
	Accepted bool `json:"accepted"`
}

type GroupMessage struct {
	ID 				 int		`json:"id"`
    GroupID          int       `json:"group_id"`
	SenderID		 string 	`json:"sender_id"`
	RecevierID		 string 	`json:"recevier_id"`
	Message 		 string      `json:"message"`
	FirstName		 string    	 `json:"first_name"`
	DateTime		 string       `json:"date_time"`
}

type TypingMessage struct {
    GroupID          int       `json:"group_id"`
	SenderID		 string 	`json:"sender_id"`
	Content 		 string     `json:"content"`
	//FirstName 		 string      `json:"first_name"`
}



type EventMessage struct {
    ID          int      `json:"id"`
	GroupID     int    `json:"group_id"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    DateTime    string   `json:"date_time"`
    Options     []Option `json:"options"`
    SenderID     string   `json:"sender_id"`
	FirstName	 string   `json:"first_name"`
	CreatedAt    string    `json:"created_at"`
	EventID      int 		`json:"event_id"`
	Day          string    	`json:"day"`
}



type Option struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
}

type EventResponseMessage struct {
	GroupID  int `json:"group_id"`
    EventID  int `json:"event_id"`
    OptionID int `json:"option_id"`
    SenderID   string `json:"sender_id"`
	FirstName  string `json:"first_name"`
}


type EventNotification struct {
	EventID      int 		`json:"event_id"`
	GroupID     int    `json:"group_id"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
	 SenderID     string   `json:"sender_id"`
	FirstName	 string   `json:"first_name"`
}


type GroupPost struct {
	GroupID       int `json:"group_id"`
	PostID        string  `json:"post_id"`
	UserID        string `json:"user_id"`
	UserNickname  string `json:"user_nickname"`
	PostTitle     string `json:"post_title"`
	Content       string `json:"content"`
	PostImage     string `json:"post_image"` // Will contain base64 string after conversion
	NumOfLikes    int    `json:"num_of_likes"`
	NumOfDislikes int    `json:"num_of_dislikes"`
	NumOfComments int    `json:"num_of_comments"`
	CreatedAt     string `json:"created_at"`
	ImageMimeType string `json:"image_mime_type"` // For content-type header
	ImageDataURL  []byte `json:"-"`               // Temporary storage, won't be sent in JSON
}

type GroupPostInteractions struct {
	Likes    int `json:"likes"`
	Dislikes int `json:"dislikes"`
	Comments int `json:"comments"`
}

// Comment struct
type GroupComment struct {
	CommentID     string `json:"comment_id"`
	PostID        string `json:"post_id"`
	UserID        string `json:"user_id"`
	UserNickname  string `json:"user_nickname"`
	CommentText   string `json:"comment_text"`
	CreatedAt     string `json:"created_at"`
	CommentImage  string `json:"comment_image"`   // Will contain base64 string after conversion
	ImageMimeType string `json:"image_mime_type"` // For content-type header
	ImageDataURL  []byte `json:"-"`               // Temporary storage, won't be sent in JSON
}


type NewGroupComment struct {
	Comment GroupComment          `json:"comment"`
	Stats   GroupPostInteractions `json:"stats"`
}

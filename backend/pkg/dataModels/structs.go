package datamodels


// Response struct to be sent back
type Response struct {

	Code     int    `json:"code"`
	Status   string `json:"status"`
	ErrorMsg string `json:"error_msg"`
	Group 	 Group	`json:"group"`
	Users	 []User `json:"users"`
	Data     interface{} `json:"data,omitempty"`
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
	 Count   				int 					`json:"count"`
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
	Count  			 int 			 `json:"count"`
	GroupName		 string			 `json:"group_name"`
	CurrentUser      string         `json:"current_user"`
}

type GroupNotifier struct {
	//ID 				 int		`json:"id"`
	SenderID		 string 	`json:"sender_id"`
	FirstName		 string    	 `json:"first_name"`
	GroupName		 string			 `json:"group_name"`
}


type TypingMessage struct {
    GroupID          int       `json:"group_id"`
	SenderID		 string 	`json:"sender_id"`
	Content 		 string     `json:"content"`
	//FirstName 		 string      `json:"first_name"`
}

type ActiveGroupMessage struct {
	Status  string `json:"status"`
	GroupID int    `json:"group_id"`
}

type ResetCountMessage struct {
	GroupID int    `json:"group_id"`
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

type UsersInvitationListMessage struct {
	GroupID		int `json:"group_id"`
	Users   []User  `json:"users"`             
}


type GroupMembersMessage struct {
	GroupID		int `json:"group_id"`
	Users   []User  `json:"users"`             
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
// UserBasicInfo struct for returning basic user information in listings
type UserBasicInfo struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}


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


type GroupPostInteractions struct {
	Likes    int `json:"likes"`
	Dislikes int `json:"dislikes"`
	Comments int `json:"comments"`
}

// Post interactions stats
type PostInteractions struct {
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

type NewComment struct {
	Comment Comment          `json:"comment"`
	Stats   PostInteractions `json:"stats"`
}

// Comment struct
type Comment struct {
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


type ExploreLists struct {
	UsersList  []User  `json:"users_list"`
	AllGroupsList []Group `json:"all_groups_list"`
	MyGroupsList []Group `json:"my_groups_list"`
	NotMyGroupsList []Group `json:"not_my_groups_list"`
}

type Profile struct {
	UserID             string `json:"id"`
	UserEmail          string `json:"email"`
	UserNickname       string `json:"nickname"`
	UserFirstName      string `json:"first_name"`
	UserLastName       string `json:"last_name"`
	UserDOB            string `json:"dob"`
	UserAvatar         string `json:"avatar"`
	UserAvatarMimeType string `json:"avatar_mime_type"`
	UserAbout          string `json:"about"`
	IsPrivate          bool   `json:"is_private"`
	UserCreatedAt      string `json:"created_at"`
	IsMyProfile        bool   `json:"is_my_profile"`
	IsFollowingMe      bool   `json:"is_following_me"`
	IsFollowingHim     bool   `json:"is_following_him"`
	IsRequestSent      bool   `json:"is_request_sent"`
	UserAvatarURL      string `json:"avatar_url"`
	NumOfFollowers     int    `json:"num_of_followers"`
	NumOfFollowing     int    `json:"num_of_following"`
	NumOfPosts         int    `json:"num_of_posts"`
}

type PrivacyUpdateRequest struct {
	IsPrivate bool `json:"is_private"`
}


type FollowRequest struct {
	RequestID string `json:"request_id"`
	UserID    string `json:"user_id"`
	UserNickname string `json:"nickname"`
}

type FollowersFollowingRequests struct {
	FollowersList []User `json:"followers_list"`
	FollowingList []User `json:"following_list"`
	RequestsList  []FollowRequest `json:"requests_list"`
}



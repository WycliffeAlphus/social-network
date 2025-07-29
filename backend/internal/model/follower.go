package model

// UserInfo represents basic user information in the response
type UserInfo struct {
	ID     string `json:"id"`
	FName  string `json:"fname"`
	LName  string `json:"lname"`
	ImgURL string `json:"imgurl"`
	Status string `json:"status,omitempty"` // only relevant for followers
}

// FollowersResponse represents the response for followers/following lists
type FollowersResponse struct {
	Users         []UserInfo `json:"users"`
	CurrentUserId string     `json:"current_user_id"`
	RequestedID   string     `json:"requested_id"`
}

type FollowRequest struct {
	FollowerID     string `json:"follower_id"`
	FollowerFname  string `json:"follower_fname"`
	FollowerLname  string `json:"follower_lname"`
	FollowerAvatar string `json:"follower_avatar"`
}

package main

// Registration info for user
type UserRegInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

//
type NewChannelInput struct {
	OwnerUID    string            `json:"ownerUID"`
	ChannelName string            `json:"channelName"`
	AccessType  ChannelAccessType `json:"accessType"`
}

// http://localhost:8080/{channel}/message
type MessageInput struct {
	ChannelName    string `json:"channelName"`
	Message        string `json:"message"`
	SenderUsername string `json:"senderUsername"`
}

// http://localhost:8080/users/{username1}

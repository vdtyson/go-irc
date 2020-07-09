package main

// http://localhost:8080/register - POST
type UserRegInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

// http://localhost:8080/channel/new - POST
type NewGroupChannelInput struct {
	OwnerUID    string            `json:"ownerUID"`
	ChannelName string            `json:"channelName"`
	AccessType  ChannelAccessType `json:"accessType"`
}

// http://localhost:8080/channels/message - POST
type NewMessageInput struct {
	ChannelName    string `json:"channelName"`
	Message        string `json:"message"`
	SenderUsername string `json:"senderUsername"`
}

// http://localhost:8080/channels/messages - GET
type AllChannelMessagesInput struct {
	UserName    string `json:"username"`
	ChannelName string `json:"channelName"`
}

// http://localhost:8080/channels/users/kick - POST TODO: Not yet implemented
type KickUserInput struct {
	ChannelName   string `json:"channelName"`
	OwnerUsername string `json:"ownerUsername"`
	UserToKick    string `json:"userToKick"`
}

// http://localhost:8080/channels/users - PUT
type AddUserToChannelInput struct {
	ChannelName   string `json:"channelName"`
	OwnerUsername string `json:"ownerUsername"`
	UserToAdd     string `json:"userToAdd"`
	PrivilegeType string `json:"privilegeType"`
}

// http://localhost:8080/users/{username}/channels - GET TODO: Not yet implemented
type AllUserChannelsInput struct {
	Username string `json:"username"`
}

/*
	PATHS without body
*/

// New DM Channel: http://localhost:8080/channels/direct/{username1}/{username2} - POST

func deleteLater() {
	//TODO: Delete later
}

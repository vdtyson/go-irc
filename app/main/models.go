package main

import (
	"time"
)

// hash tag for group chats ex: #test
// direct message !@testuser@testuser2

// users -> [uid]
type User struct {
	IsAdmin  bool   `json:"isAdmin"`
	IsBanned bool   `json:"isBanned"`
	Username string `json:"username"`
}

// Channel user is member of
// user -> [uid] -> userchannels -> [#channelname]
type UserPrivilegeType string // privileges a user has in a channel
const (
	OWNER  = UserPrivilegeType("OWNER")  // privileges: define channel access type, write,kick member,add member, change other members access type
	MOD    = UserPrivilegeType("MOD")    // privileges: kick member,write
	MEMBER = UserPrivilegeType("MEMBER") // privileges: write
)

type UserChannel struct {
	PrivilegeType UserPrivilegeType `json:"privilegeType"`
}

// channels -> [#channelname]
type ChannelAccessType string // defines how users can join; channel owner can set this
const (
	CLOSED = ChannelAccessType("DIRECT_MESSAGE") // direct message
	INVITE = ChannelAccessType("INVITE")         // channel is invite only
	SECRET = ChannelAccessType("SECRET")         // can join through secret key
	OPEN   = ChannelAccessType("PUBLIC")         // can join through channel name
)

type Channel struct {
	AccessType   ChannelAccessType `json:"accessType"`
	LastMessages []*Message        `json:"lastMessages"` // messageId as base map key; child map contains message attributes (message, senderId, timeSpent)
}

// usernames -> [username]
// doc has username as key to make unique; has a field that points to the uid
type Username struct {
	Uid string `json:"uid"`
}

// channelChats -> [channelName] -> messages -> [messageId] // channelChats holds messages of a specific channel based off channelName
type Message struct {
	Value          string    `json:"message"`
	SenderId       string    `json:"senderId"`
	SenderUsername string    `json:"senderUsername"`
	TimeSent       time.Time `json:"timeSent"`
}

// Registration info for user
type UserRegInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

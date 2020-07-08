package main

import (
	"time"
)

// hash tag for group chats ex: #test
// direct message !@testuser@testuser2
type DocModel interface {
	Map() map[string]interface{}
	ID() string
}

// users -> [uid]
type User struct {
	uid      string
	isAdmin  bool
	isBanned bool
	username string
}

func (u *User) Map() map[string]interface{} {
	return map[string]interface{}{
		"isAdmin":  u.isAdmin,
		"isBanned": u.isBanned,
		"username": u.username,
	}
}
func (u *User) ID() string {
	return u.uid
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
	channelName   string
	privilegeType UserPrivilegeType
}

func (u *UserChannel) Map() map[string]interface{} {
	return map[string]interface{}{
		"channelName":   u.channelName,
		"privilegeType": string(u.privilegeType),
	}
}
func (u *UserChannel) ID() string {
	return u.channelName
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
	channelName  string
	accessType   ChannelAccessType
	lastMessages map[string]map[string]interface{} // messageId as base map key; child map contains message attributes (message, senderId, timeSpent)
}

func (c *Channel) Map() map[string]interface{} {
	return map[string]interface{}{
		"accessType":   string(c.accessType),
		"lastMessages": c.lastMessages, // messageId as base map key; child map contains message attributes (message, senderId, timeSpent)
	}
}
func (c *Channel) ID() string {
	return c.channelName
}

// usernames -> [username]
// doc has username as key to make unique; has a field that points to the uid
type Username struct {
	value string
	uid   string
}

func (u *Username) Map() map[string]interface{} {
	return map[string]interface{}{"uid": u.uid}
}
func (u *Username) ID() string {
	return u.value
}

// channelChats -> [channelName] -> messages -> [messageId] // channelChats holds messages of a specific channel based off channelName
type Message struct {
	messageId      string
	value          string
	senderId       string
	senderUsername string
	timeSent       time.Time
}

func (m *Message) Map() map[string]interface{} {
	return map[string]interface{}{
		"message":        m.value,
		"senderId":       m.senderId,
		"senderUsername": m.senderUsername,
		"timeSent":       m.timeSent,
	}
}

func (m *Message) ID() string {
	return m.messageId
}

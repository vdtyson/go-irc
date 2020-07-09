package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	CHANNELS_PATH = "channels" // channels -> [channelName]
	MEMBERS_PATH  = "members"  // channels -> [channelName] -> members

	USERS_PATH         = "users"        // users -> [uid]
	USER_CHANNELS_PATH = "userChannels" // users -> [uid] -> userChannels

	USERNAMES_PATH = "usernames" // usernames -> [username]

	CHANNEL_CHATS_PATH = "channelChats" // channelChats ->
	MESSAGES_PATH      = "messages"     // channelChats -> [channelName] -> messages
)

type Repository interface{}
type FirestoreRepository interface {
	Repository
}

type AuthRepository struct {
	authClient      *auth.Client
	firestoreClient *firestore.Client
}

type UsernameExistsError struct {
	username string
}

func (u *UsernameExistsError) Error() string {
	return fmt.Sprintf("Username %s already exists.", u)
}

/*func (a *AuthRepository) SignIn(email, password string) (*auth.UserRecord, error) {
	a.authClient.S
}*/
// true if account created successfully

func (a *AuthRepository) RegisterUser(ctx context.Context, userRegInfo UserRegInfo) (*auth.UserRecord, error) {

	fmt.Println("started method registerUser()")
	userNameRef := a.firestoreClient.Collection(USERNAMES_PATH).Doc(userRegInfo.Username)

	docSnapshot, err := userNameRef.Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}

	if docSnapshot.Exists() {
		return nil, &UsernameExistsError{userRegInfo.Username}
	}

	params := (&auth.UserToCreate{}).
		Email(userRegInfo.Email).
		Password(userRegInfo.Password)

	authUser, err := a.authClient.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	fmt.Printf("authUser: %+v\n", authUser)

	user := User{IsAdmin: userRegInfo.IsAdmin, IsBanned: false, Username: userRegInfo.Username}
	username := Username{Uid: authUser.UID}
	userChannel := UserChannel{PrivilegeType: MEMBER}

	userRef := a.firestoreClient.Collection(USERS_PATH).Doc(authUser.UID)
	channelRef := a.firestoreClient.Collection(CHANNELS_PATH).Doc("#main").Collection(MEMBERS_PATH).Doc(userRegInfo.Username)

	batch := a.firestoreClient.Batch()

	_, err = batch.
		Set(userNameRef, &username).
		Set(userRef, &user).
		Set(userRef.Collection(USER_CHANNELS_PATH).Doc("#main"), &userChannel).
		Set(userNameRef, &username).
		Create(channelRef, map[string]interface{}{}).
		Commit(ctx)

	if err != nil {
		return nil, err
	}
	return authUser, nil
}

type ChannelRepository struct {
	fsClient *firestore.Client
}

func (c *ChannelRepository) CreateDirectMessageChannel(ctx context.Context, username1, username2 string) error {

	channelName := fmt.Sprintf("!@%s@%s", username1, username2)
	messages := make([]*Message, 0)

	usernameRef := c.fsClient.Collection(USERNAMES_PATH)

	var u1 Username
	username1DocSnapshot, err := usernameRef.Doc(username1).Get(ctx)
	if err != nil {
		return err
	}
	err = username1DocSnapshot.DataTo(&u1)
	if err != nil {
		return err
	}

	var u2 Username
	username2DocSnapshot, err := usernameRef.Doc(username2).Get(ctx)
	if err != nil {
		return err
	}
	err = username2DocSnapshot.DataTo(&u2)
	if err != nil {
		return err
	}

	channel := Channel{AccessType: CLOSED, LastMessages: messages}
	userChannel := UserChannel{PrivilegeType: OWNER}

	channelRef := c.fsClient.Collection(CHANNELS_PATH).Doc(channelName)
	userRef := c.fsClient.Collection(USERS_PATH)

	batch := c.fsClient.Batch()

	_, err = batch.
		Set(channelRef, &channel).
		Create(channelRef.Collection(MEMBERS_PATH).Doc(username1), map[string]interface{}{}).
		Create(channelRef.Collection(MEMBERS_PATH).Doc(username2), map[string]interface{}{}).
		Set(userRef.Doc(u1.Uid).Collection(USER_CHANNELS_PATH).Doc(channelName), &userChannel).
		Set(userRef.Doc(u2.Uid).Collection(USER_CHANNELS_PATH).Doc(channelName), &userChannel).
		Commit(ctx)

	return err
}

func (c *ChannelRepository) CreateGroupChannel(ctx context.Context, channelInput NewChannelInput) error {
	channelName := fmt.Sprintf("#%s", channelInput.ChannelName)
	messages := make([]*Message, 0)

	channel := Channel{AccessType: channelInput.AccessType, LastMessages: messages}
	userChannel := UserChannel{PrivilegeType: OWNER}

	channelRef := c.fsClient.Collection(CHANNELS_PATH).Doc(channelName)
	userChannelRef := c.fsClient.Collection(USERS_PATH).Doc(channelInput.OwnerUID).Collection(USER_CHANNELS_PATH).Doc(channelName)

	batch := c.fsClient.Batch()

	_, err := batch.
		Set(channelRef, &channel).
		Create(channelRef.Collection(MEMBERS_PATH).Doc(channelInput.OwnerUID), map[string]interface{}{}).
		Set(userChannelRef, &userChannel).
		Commit(ctx)

	return err
}

// http://localhost:8080/channels/{channelName}/messages/all
func (c *ChannelRepository) GetAllChannelMessages(ctx context.Context, input ChannelNameInput) ([]*Message, error) {
	messageRefs, err := c.fsClient.Collection(CHANNEL_CHATS_PATH).Doc(input.ChannelName).Collection(MESSAGES_PATH).DocumentRefs(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var messages []*Message
	for _, messageRef := range messageRefs {
		messageSnapshot, err := messageRef.Get(ctx)
		if err != nil {
			return nil, err
		}
		var message Message
		err = messageSnapshot.DataTo(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// http://localhost:8080/{channel}/message
func (c *ChannelRepository) NewMessage(ctx context.Context, messageInput MessageInput) error {
	chanChatsRef := c.fsClient.Collection(CHANNEL_CHATS_PATH).Doc(messageInput.ChannelName).Collection(MESSAGES_PATH)
	chanRef := c.fsClient.Collection(CHANNELS_PATH).Doc(messageInput.ChannelName)
	chanDocSnapshot, err := chanRef.Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}

	if !chanDocSnapshot.Exists() {
		return fmt.Errorf("channel does not exist")
	}

	memberDocSnapshot, err := chanRef.Collection(MEMBERS_PATH).Doc(messageInput.SenderUsername).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}

	if !memberDocSnapshot.Exists() {
		return fmt.Errorf("user does not have access to this channel")
	}

	var channel Channel
	err = chanDocSnapshot.DataTo(&channel)
	if err != nil {
		return err
	}

	newMessage := Message{SenderUsername: messageInput.SenderUsername, SenderMessage: messageInput.Message, TimeSent: time.Now()}

	if len(channel.LastMessages) >= 5 {
		newRecents := []*Message{
			channel.LastMessages[1],
			channel.LastMessages[2],
			channel.LastMessages[3],
			channel.LastMessages[4],
			&newMessage,
		}

		channel.LastMessages = newRecents
	} else {
		channel.LastMessages = append(channel.LastMessages, &newMessage)
	}

	_, err = chanRef.Set(ctx, &channel)
	if err != nil {
		return err
	}
	_, _, err = chanChatsRef.Add(ctx, &newMessage)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChannelRepository) KickUser(ctx context.Context, input KickUserInput) error {
	// check if kicker is admin of channel
	// check if user exists in channel
	// delete channel from userChannel and member from [channelName] -> member
	channelRef := c.fsClient.Collection(CHANNELS_PATH).Doc(input.ChannelName)
	_, err := channelRef.Get(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChannelRepository) AddUser(ctx context.Context, input AddUserInput) error {
	var ownerUsername Username
	var userToAddUsername Username

	channelRef := c.fsClient.Collection(CHANNELS_PATH).Doc(input.ChannelName)
	_, err := channelRef.Get(ctx)
	if err != nil {
		return err
	}

	ownerUserNameSnapshot, err := c.fsClient.Collection(USERNAMES_PATH).Doc(input.OwnerUsername).Get(ctx)
	if err != nil {
		return err
	}
	err = ownerUserNameSnapshot.DataTo(&ownerUsername)
	if err != nil {
		return err
	}

	userToAddUsernameSnapshot, err := c.fsClient.Collection(USERNAMES_PATH).Doc(input.UserToAdd).Get(ctx)
	if err != nil {
		return err
	}
	err = userToAddUsernameSnapshot.DataTo(&userToAddUsername)
	if err != nil {
		return err
	}

	var ownerUserChannel UserChannel
	ownerUserChannelSnapshot, err := c.fsClient.Collection(USERS_PATH).Doc(ownerUsername.Uid).Collection(USER_CHANNELS_PATH).Doc(input.ChannelName).Get(ctx)
	if err != nil {
		return err
	}
	err = ownerUserChannelSnapshot.DataTo(&ownerUserChannel)
	if err != nil {
		return err
	}

	if ownerUserChannel.PrivilegeType != OWNER {
		return fmt.Errorf("user is not an owner of this channel")
	}

	userChannel := UserChannel{PrivilegeType: UserPrivilegeType(input.PrivilegeType)}
	batch := c.fsClient.Batch()
	_, err = batch.
		Set(c.fsClient.Collection(USERS_PATH).Doc(userToAddUsername.Uid).Collection(USER_CHANNELS_PATH).Doc(input.ChannelName), &userChannel).
		Create(channelRef.Collection(MEMBERS_PATH).Doc(input.UserToAdd), map[string]interface{}{}).
		Commit(ctx)

	return err
}

// TODO: Kick user
// TODO: Ban User
// TODO: Get all user channels by username
// TODO: Join channel

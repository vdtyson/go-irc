package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	channelRef := a.firestoreClient.Collection(CHANNELS_PATH).Doc("#main").Collection(MEMBERS_PATH).Doc(authUser.UID)

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

func (a *AuthRepository) GetUser() {

}

func (a *AuthRepository) Login() {

}

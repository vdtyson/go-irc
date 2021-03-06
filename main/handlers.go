package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	JSON_PREFIX = ""
	JSON_INDENT = "    "
)

// GetNewestMessage godoc
// @Summary Get newest message from channel
// @Param channelMessagesInput body AllChannelMessagesInput true "Channel messages input"
// @Accept json
// @Router /channels/messages/newest [post]
// @Tags channels
func GetNewestMessageHandler(w http.ResponseWriter, r *http.Request) {
	var messageInput AllChannelMessagesInput
	err := json.NewDecoder(r.Body).Decode(&messageInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	newestMessage, err := appInstance.channelRepo.GetNewestMessage(r.Context(), messageInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data, err := json.MarshalIndent(newestMessage, JSON_PREFIX, JSON_INDENT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Write(data)
}

// RegisterUser godoc
// @Summary Register a new user
// @Param userRegistrationInput body UserRegInput true "User registration info"
// @Accept json
// @Router /register [post]
// @Tags users
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var userRegInfo UserRegInput

	err := json.NewDecoder(r.Body).Decode(&userRegInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = appInstance.authRepo.RegisterUser(r.Context(), userRegInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// CreateGroupChannel godoc
// @Summary Create a group channel
// @Param newGroupChannelInput body NewGroupChannelInput true "Input for new group channel"
// @Accept json
// @Router /channels/new [post]
// @Tags channels
func CreateGroupChannelHandler(w http.ResponseWriter, r *http.Request) {
	var newChannelInput NewGroupChannelInput

	err := json.NewDecoder(r.Body).Decode(&newChannelInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = appInstance.channelRepo.CreateGroupChannel(r.Context(), newChannelInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// NewMessage godoc
// @Summary Send new message to channel
// @Param newMessageInput body NewMessageInput true "Input for new message"
// @Tags users
// @Router /channels/message [post]
func NewMessageHandler(w http.ResponseWriter, r *http.Request) {
	var messageInput NewMessageInput

	err := json.NewDecoder(r.Body).Decode(&messageInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = appInstance.channelRepo.NewMessage(r.Context(), messageInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// http://localhost:8080/channels/direct/{username1}/{username2}
func CreateDMChannelHandler(w http.ResponseWriter, r *http.Request) {
	username1, username2 := mux.Vars(r)["username1"], mux.Vars(r)["username2"]

	err := appInstance.channelRepo.CreateDirectMessageChannel(r.Context(), username1, username2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetAllMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var channelNameInput AllChannelMessagesInput
	var body []byte
	r.Body.Read(body)
	fmt.Println("body:", string(body))

	err := json.NewDecoder(r.Body).Decode(&channelNameInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Printf("channel name input: %+v\n", channelNameInput)

	messages, err := appInstance.channelRepo.GetAllChannelMessages(r.Context(), channelNameInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data, err := json.MarshalIndent(messages, JSON_PREFIX, JSON_INDENT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	var addUserInput AddUserToChannelInput

	err := json.NewDecoder(r.Body).Decode(&addUserInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !IsPrivilegeType(addUserInput.PrivilegeType) {
		http.Error(w, fmt.Errorf("value %s is an invalid privelegeType", addUserInput.PrivilegeType).Error(), http.StatusBadRequest)
		return
	}

	err = appInstance.channelRepo.AddUser(r.Context(), addUserInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUserChannelsHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userChannels, err := appInstance.userRepo.GetAllUserChannels(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data, err := json.MarshalIndent(userChannels, JSON_PREFIX, JSON_INDENT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func BanUserHandler(w http.ResponseWriter, r *http.Request) {
	var banUserInput BanUserInput

	err := json.NewDecoder(r.Body).Decode(&banUserInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = appInstance.adminRepo.BanUser(r.Context(), banUserInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("User %s was banned.", banUserInput.UserToBanUsername)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

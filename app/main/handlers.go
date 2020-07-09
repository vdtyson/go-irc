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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var userRegInfo UserRegInfo

	err := json.NewDecoder(r.Body).Decode(&userRegInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = appInstance.authRepo.RegisterUser(r.Context(), userRegInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func CreateGroupChannelHandler(w http.ResponseWriter, r *http.Request) {
	var newChannelInput NewChannelInput

	err := json.NewDecoder(r.Body).Decode(&newChannelInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = appInstance.channelRepo.CreateGroupChannel(r.Context(), newChannelInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func NewMessageHandler(w http.ResponseWriter, r *http.Request) {
	var messageInput MessageInput

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
	err := json.NewDecoder(r.Body).Decode(&channelNameInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

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
	var addUserInput AddUserInput

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

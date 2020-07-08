package main

import (
	"encoding/json"
	"net/http"
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

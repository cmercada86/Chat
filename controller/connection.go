package controller

import (
	"Chat/auth"
	"Chat/cache"
	"Chat/model"
	"Chat/repository"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Connect(w http.ResponseWriter, r *http.Request) {

	x, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	state := model.UUID(vars["state"])

	if !cache.Contains(state) {
		errorHandler(w, r, 403, "Invalid State")
		return
	}
	//get auth code	string
	authCode := string(x)

	user, isPlusUser, err := auth.ExchangeAuthCodeForUser(authCode)
	if err != nil {
		log.Println("Error getting Auth code from Google: ", err)
		errorHandler(w, r, 403, "Could not get user id from Google")
		return
	}
	if !isPlusUser {
		//SEND UNAUTH
		errorHandler(w, r, 403, "Must be a Google+ user!")
		return
	}

	log.Println(user.Name, " Logged in!")
	cache.Set(state, user)

	repository.AddOrUpdateUserInfo(user)

	json.NewEncoder(w).Encode(user)

}

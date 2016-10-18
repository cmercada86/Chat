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
		log.Println("invalid state!")
		return
	}
	//get auth code	string
	authCode := string(x)

	user, isPlusUser, err := auth.ExchangeAuthCodeForUser(authCode)
	if err != nil {
		log.Println("Error getting Auth code from Google: ", err)
	}
	if !isPlusUser {
		//SEND UNAUTH
	}

	log.Println(user.Name, " Logged in!")
	cache.Set(state, user)

	repository.AddOrUpdateUserInfo(user)

	json.NewEncoder(w).Encode(user)

}

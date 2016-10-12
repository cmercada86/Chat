package controller

import (
	"Chat/auth"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func connect(w http.ResponseWriter, r *http.Request) {
	log.Println("HERE")
	x, _ := ioutil.ReadAll(r.Body)

	//get auth code	string
	authCode := string(x)

	user, err := auth.ExchangeAuthCodeForUser(authCode)
	if err != nil {

	}

	json.NewEncoder(w).Encode(user)

}

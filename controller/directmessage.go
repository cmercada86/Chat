package controller

import (
	"Chat/auth"
	"Chat/model"
	"Chat/repository"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SendDirectMessage(w http.ResponseWriter, r *http.Request) {
	x, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	state := model.UUID(vars["state"])
	receiverID := vars["receiver_id"]

	user, isAuth := auth.CheckAuth(state)
	if !isAuth {
		errorHandler(w, r, 404, "Not Authorized")
		return
	}

	message := string(x)
	dm := model.DirectMessage{
		SenderID:   user.ID,
		ReceiverID: receiverID,
		Message:    message,
		Seen:       false,
	}

	//dont want to send message to yourself!
	if user.ID == receiverID {
		return
	}
	log.Println(dm)

	repository.InsertDirectMessage(dm)
}

func ReceiveDirectMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	state := model.UUID(vars["state"])

	user, isAuth := auth.CheckAuth(state)
	if !isAuth {
		errorHandler(w, r, 404, "Not Authorized")
		return
	}

	dms, err := repository.GetDirectMessages(user.ID)
	if err != nil {
		errorHandler(w, r, 500, "")
		return
	}

	json.NewEncoder(w).Encode(dms)

}

package main

import (
	"Chat/auth"
	"Chat/config"
	"Chat/route"
	"log"
	"net/http"
)

func main() {

	config.ReadConfFile("")

	con := config.GetConfig()

	auth.SetOAuth2Config(con.GogClientID, con.GogClientSecret)

	router := route.NewRouter()
	//http.HandleFunc("/connect", connect)
	//http.HandleFunc("/", index)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Listen failed: ", err)
	}

}

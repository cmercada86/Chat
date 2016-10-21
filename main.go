package main

import (
	"Chat/auth"
	"Chat/cache"
	"Chat/config"
	"Chat/repository"
	"Chat/route"
	"log"
	"net/http"
)

func main() {

	config.ReadConfFile("config.json")

	con := config.GetConfig()

	cache.ConnectToCache(con.RedisURL)
	defer cache.Close()

	repository.NewRepository(con.PostgresUser, con.PostgresPass, con.PostgresHost)
	repository.NewDBtracker(con.PostgresUser, con.PostgresPass, con.PostgresHost)

	repository.SetSearchUrl(con.SearchURL)

	go repository.Listen()

	auth.SetOAuth2Config(con.GogClientID, con.GogClientSecret)

	router := route.NewRouter()
	//http.HandleFunc("/connect", connect)
	//http.HandleFunc("/", index)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Listen failed: ", err)
	}

}

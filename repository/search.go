package repository

import (
	"Chat/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

var url string

type search struct {
	Room         string `json:"room"`
	SearchString string `json:"search"`
}

type resultStruct struct {
	Uid       string    `json:"uid"`
	User_ID   string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Room      string    `json:"room"`
	Message   string    `json:"message"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
}

func SetSearchUrl(searchUrl string) {
	url = searchUrl
}

func SearchChat(room string, searchString string) ([]model.Chat, error) {
	var chats []model.Chat

	conn, _ := net.Dial("tcp", url)
	defer conn.Close()

	fmt.Fprintf(conn, model.ObjectToJsonString(search{
		Room:         room,
		SearchString: searchString,
	}))

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return chats, err
	}

	var resList []resultStruct

	err = json.Unmarshal(result, &resList)
	for _, res := range resList {
		chats = append(chats, model.Chat{
			Uid:       res.Uid,
			Timestamp: res.Timestamp,
			User: model.User{
				ID:      res.User_ID,
				Name:    res.Name,
				Picture: res.Picture,
			},
			Room:    res.Room,
			Message: res.Message,
		})
	}

	return chats, err
}

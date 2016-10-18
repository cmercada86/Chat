package model

import (
	"encoding/json"
	"time"
)

type Chat struct {
	Uid       []uint8   `json:"uuid"`
	Timestamp time.Time `json:"timestamp"`
	User      User      `json:"user"`
	Room      string    `json:"room"`
	Message   string    `json:"message"`
}

type DirectMessage struct {
	Uid        UUID
	Timestamp  time.Time `json:"timestamp"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Message    string    `json:"message"`
	Seen       bool      `json:"seen"`
}

func ObjectToJsonString(ob interface{}) string {
	b, _ := json.Marshal(ob)

	return string(b)
}

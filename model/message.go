package model

import "time"

type Chat struct {
	Uid       UUID
	Timestamp time.Time
	UserID    string
	Room      string
	Message   string
}

type DirectMessage struct {
	Uid        UUID
	Timestamp  time.Time `json:"timestamp"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Message    string    `json:"message"`
	Seen       bool      `json:"seen"`
}

func ChatToJsonString(chat Chat) string {

}
func DirectMessageToJsonSring(dm DirectMessage) string {

}

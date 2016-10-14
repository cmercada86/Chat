package model

import "time"

type Chat struct {
	Uid       UUID
	Timestamp time.Time
	UserID    string
	Room      string
	Message   string
}

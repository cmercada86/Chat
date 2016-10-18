package model

type Listener struct {
	UserID      string
	CurrentRoom string
	ChatChannel chan Chat
	DMchannel   chan DirectMessage
	UserChannel chan User
}

func NewListener(userID string, room string) *Listener {
	return &Listener{
		UserID:      userID,
		CurrentRoom: room,
		ChatChannel: make(chan Chat, 5),
		DMchannel:   make(chan DirectMessage, 5),
		UserChannel: make(chan User, 5),
	}
}

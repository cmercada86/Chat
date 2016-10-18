package model

type Listener struct {
	User        User
	CurrentRoom string
	ChatChannel chan Chat
	DMchannel   chan DirectMessage
	UserChannel chan User
	RoomChannel chan string
}

func NewListener(user User, room string) *Listener {
	return &Listener{
		User:        user,
		CurrentRoom: room,
		ChatChannel: make(chan Chat, 5),
		DMchannel:   make(chan DirectMessage, 5),
		UserChannel: make(chan User, 5),
		RoomChannel: make(chan string)
	}
}

func (listener *Listener) Close() {
	close(listener.ChatChannel)
	close(listener.DMchannel)
	close(listener.UserChannel)
	close(listener.RoomChannel)
}

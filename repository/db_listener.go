package repository

import (
	"Chat/model"
	"encoding/json"
	"log"
	"time"

	"github.com/streamrail/concurrent-map"

	"github.com/lib/pq"
)

var db_listener *pq.Listener
var stop chan struct{}
var listeners cmap.ConcurrentMap

type notification struct {
	Table string          `json:"table"`
	Type  string          `json:"type"`
	Data  json.RawMessage `json:"data"`
}

type Time struct {
	time.Time
}

//for unmarshalling correct dt from postgres
type chatWrapped struct {
	Uid       []uint8 `json:"uuid"`
	Timestamp Time    `json:"timestamp"`
	UserID    string  `json:"user_id"`
	Room      string  `json:"room"`
	Message   string  `json:"message"`
}
type dmWrapped struct {
	Uid        model.UUID
	Timestamp  Time   `json:"timestamp"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Message    string `json:"message"`
	Seen       bool   `json:"seen"`
}

var format = "2006-01-02T15:04:05.999999"

func NewDBtracker(user string, pass string, host string) {
	connectionString := getConnectionString(user, pass, host)

	db_listener = pq.NewListener(connectionString, 10*time.Second, time.Minute, logIssue)
	listeners = cmap.New()

	stop = make(chan struct{})

}

func Listen() {

	rooms, _ := GetRoomNames()

	if err := db_listener.Listen("table_update"); err != nil {
		log.Fatalln(err)
	}

	for {
		select {

		case n := <-db_listener.Notify:

			var notify notification

			json.Unmarshal([]byte(n.Extra), &notify)

			switch notify.Table {
			case "chat_table":
				var chatWrap chatWrapped
				if err := json.Unmarshal(notify.Data, &chatWrap); err != nil {
					log.Println("error unmarshalling chat: ", err)
				}
				chat := model.Chat{
					Uid:       chatWrap.Uid,
					Timestamp: chatWrap.Timestamp.Time,
					Room:      chatWrap.Room,
					Message:   chatWrap.Message,
				}

				user, err := GetUserFromID(chatWrap.UserID)
				if err != nil {
					log.Println("Error getting user: ", err)
				} else {
					chat.User = user
				}

				for listener := range listeners.Iter() {
					listener.Val.(*model.Listener).ChatChannel <- chat
				}

				if !stringsContain(rooms, chat.Room) {
					rooms = append(rooms, chat.Room)

					for listener := range listeners.Iter() {
						listener.Val.(*model.Listener).RoomChannel <- chat.Room
					}
				}

			case "dm_table":
				var dmWrap dmWrapped
				if err := json.Unmarshal(notify.Data, &dmWrap); err != nil {
					log.Println("error unmarshalling chat: ", err)
				}
				dm := model.DirectMessage{
					Uid:       dmWrap.Uid,
					Timestamp: dmWrap.Timestamp.Time,
					Message:   dmWrap.Message,
					Seen:	dmWrap.Seen,
				}
				

				sender, err := GetUserFromID(dm.SenderID)
				if err != nil {
					log.Println("Error getting user: ", err)
				} else {
					dm.Sender = sender
				}
				receiver, err := GetUserFromID(dm.ReceiverID)
				if err != nil {
					log.Println("Error getting user: ", err)
				} else {
					dm.Receiver = receiver
				}
				for listener := range listeners.Iter() {
					if listener.Val.(*model.Listener).User.ID==receiver.ID{
						listener.Val.(*model.Listener).DMchannel <- dm
					}
				}
			}
		case <-stop:
			db_listener.Close()
			return
		}

	}

}

func GetCurrentListeners() ([]model.User, error) {
	var users []model.User
	for listener := range listeners.Iter() {
		user := listener.Val.(*model.Listener).User
		users = append(users, user)
	}
	return users, nil

	//return GetUsersFromID(listeners.Keys())
}

func AddListener(newListener *model.Listener) {
	for listener := range listeners.Iter() {
		listener.Val.(*model.Listener).UserChannel <- newListener.User
	}
	listeners.Set(newListener.User.ID, newListener)

}
func RemoveListener(userID string) {
	listener, ok := listeners.Get(userID)
	if ok {
		listener.(*model.Listener).Close()
		listeners.Remove(userID)
	}
}

func StopTracking() {
	close(stop)
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	b = b[1 : len(b)-1]
	t.Time, err = time.Parse(format, string(b))
	return
}

func stringsContain(sList []string, obj string) bool {
	for _, s := range sList {
		if s == obj {
			return true
		}
	}
	return false
}

func logIssue(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

package repository

import (
	"Chat/model"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type MessageStatus int

const (
	New MessageStatus = iota
	Seen
	All
	StatusError
)

const maxConnections = 20

var db *sql.DB

func NewRepository(user string, pass string, host string) {
	connectionString := getConnectionString(user, pass, host)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}
	db.SetMaxOpenConns(maxConnections)
	if err := db.Ping(); err != nil {
		log.Println("Error ping: ", err)
	}
}

func AddOrUpdateUserInfo(user model.User) {
	query := `
		INSERT INTO user_table (id, name, firstname,lastname,date_added, picture, locale )
		VALUES ('$ID','$NAME','$FIRSTNAME','$LASTNAME',NOW(),'$PICTURE','$LOCALE')
		ON CONFLICT(id)  DO UPDATE
		SET name='$NAME',firstname='$FIRSTNAME',lastname='$LASTNAME',
			picture='$PICTURE', locale='$LOCALE'
		;`
	query = strings.Replace(query, "$ID", user.ID, -1)
	query = strings.Replace(query, "$NAME", user.Name, -1)
	query = strings.Replace(query, "$FIRSTNAME", user.FirstName, -1)
	query = strings.Replace(query, "$LASTNAME", user.LastName, -1)
	query = strings.Replace(query, "$PICTURE", user.Picture, -1)
	query = strings.Replace(query, "$LOCALE", user.Locale, -1)

	if err := db.QueryRow(query).Scan(); err != nil && err != sql.ErrNoRows {
		log.Println("Error adding or updating user: ", err, query)
	}

}

func AddChatMessage(room string, userID string, message string) {
	query :=
		`
		INSERT INTO chat_table (uid, user_id, timestamp,room, message)
		VALUES (uuid_generate_v4(),$1,NOW(),$2,$3)
		;`

	if err := db.QueryRow(query, userID, room, message).Scan(); err != nil && err != sql.ErrNoRows {
		log.Println("Error adding or updating chat: ", err, query)
	}
	//log.Println("Room: ", room, " Message: ", message)
}

func GetChatMessages(room string) ([]model.Chat, error) {

	var chats []model.Chat
	query := fmt.Sprintf(`SELECT * FROM chat_table
	JOIN user_table ON (chat_table.user_id = user_table.id)
	WHERE room='%s';`, room)

	rows, err := db.Query(query)
	if err != nil && err != sql.ErrNoRows {
		//log.Println("Error querying chat_table: ", err)
		return chats, err
	}

	for rows.Next() {
		chat, err := readChat(rows)
		if err != nil {

		} else {
			chats = append(chats, chat)
		}
	}

	return chats, nil

}

func GetRoomNames() ([]string, error) {
	var rooms []string

	rows, err := db.Query("SELECT DISTINCT room from chat_table;")
	if err != nil && err != sql.ErrNoRows {
		//log.Println("Error querying chat_table: ", err)
		return rooms, err
	}
	for rows.Next() {
		var room string
		rows.Scan(&room)

		rooms = append(rooms, room)

	}
	return rooms, nil
}

func GetUserFromID(userID string) (model.User, error) {

	row := db.QueryRow("SELECT * from user_table WHERE id='" + userID + "';")

	return readUserRow(row)

}

func GetUsersFromID(userIDs []string) ([]model.User, error) {
	var users []model.User

	query := "SELECT * FROM user_table WHERE "

	for i, userID := range userIDs {
		if i > 0 {
			query += " OR "
		}
		query += fmt.Sprintf("id='%s'", userID)
	}

	rows, err := db.Query(query + ";")
	if err != nil && err != sql.ErrNoRows {
		//log.Println("Error querying chat_table: ", err)
		return users, err
	}

	for rows.Next() {
		user, err := readUser(rows)
		if err != nil {
			//
		} else {
			users = append(users, user)
		}

	}
	return users, nil
}

func InsertDirectMessage(senderID string, receiverID string, message string) {
	query :=
		`
		INSERT INTO dm_table (uid, sender_id, timestamp,receiver_id, message,seen)
		VALUES (uuid_generate_v4(),$1,NOW(),$2,$3,false)
		;`

	if err := db.QueryRow(query, senderID, receiverID, message).Scan(); err != nil && err != sql.ErrNoRows {
		log.Println("Error adding or updating chat: ", err, query)
	}
}

func GetDirectMessages(receiverID string, status MessageStatus) ([]model.DirectMessage, error) {
	var dms []model.DirectMessage
	query := fmt.Sprintf(`SELECT * FROM dm_table as dms
	 JOIN user_table as send_user ON (dms.sender_id = send_user.id)
	 JOIN user_table as receive_user ON (dms.receiver_id = receive_user.id)
	WHERE receiver_id='%s'`, receiverID)

	if status == New || status == Seen {
		query = fmt.Sprintf("%s AND seen=%v;", query, status == Seen)
	} else {
		query += ";"
	}

	rows, err := db.Query(query)
	if err != nil && err != sql.ErrNoRows {
		//log.Println("Error querying chat_table: ", err)
		return dms, err
	}

	for rows.Next() {
		dm, err := readDM(rows)
		if err != nil {

		} else {
			dms = append(dms, dm)
		}
	}

	return dms, nil
}

func SetDirectMessageSeen(dmID string) {
	query := fmt.Sprintf(
		`
		UPDATE dm_table 
		SET seen = true
		WHERE uid='%s'
		;`, dmID)

	if err := db.QueryRow(query).Scan(); err != nil && err != sql.ErrNoRows {
		log.Println("Error adding or updating dm: ", err, query)
	}
}

func GetMessageStatusFromString(status string) MessageStatus {
	switch status {
	case "New":
		return New
	case "Seen":
		return Seen
	case "All":
		return All
	default:
		return StatusError
	}

}

func Close() {
	db.Close()
}

func getConnectionString(user string, pass string, host string) string {
	return fmt.Sprintf("user=%s password=%s dbname=chat host=%s sslmode=disable ",
		user,
		pass,
		host)
}

func readChat(rows *sql.Rows) (model.Chat, error) {
	var uuid []byte
	var chat model.Chat
	var userID string
	var user model.User
	var tempDate time.Time

	if err := rows.Scan(&uuid, &userID, &chat.Timestamp, &chat.Room,
		&chat.Message, &user.ID, &user.Name, &user.FirstName, &user.LastName,
		&tempDate, &user.Picture, &user.Locale); err != nil {
		log.Println("Error reading chat: ", err)
		return chat, err
	}
	chat.Uid = string(uuid)

	//chat.User, err = readUser(rows)

	chat.User = user

	return chat, nil

}

func readDM(rows *sql.Rows) (model.DirectMessage, error) {
	var uuid []byte
	var dm model.DirectMessage
	var sendID string
	var receiveID string
	var sender model.User
	var receiver model.User
	var tempDate time.Time
	if err := rows.Scan(&uuid, &sendID, &dm.Timestamp, &receiveID,
		&dm.Message, &dm.Seen, &sender.ID, &sender.Name, &sender.FirstName, &sender.LastName,
		&tempDate, &sender.Picture, &sender.Locale, &receiver.ID, &receiver.Name,
		&receiver.FirstName, &receiver.LastName, &tempDate, &receiver.Picture, &receiver.Locale,
	); err != nil {
		log.Println("Error reading chat: ", err)
		return dm, err
	}
	dm.Uid = string(uuid)
	dm.Sender = sender
	dm.Receiver = receiver

	return dm, nil

}

func readUser(rows *sql.Rows) (model.User, error) {
	var user model.User
	var date time.Time
	if err := rows.Scan(&user.ID, &user.Name, &user.FirstName, &user.LastName,
		&date, &user.Picture, &user.Locale); err != nil {
		log.Println("Error reading user: ", err)
		return user, err
	}

	return user, nil
}
func readUserRow(row *sql.Row) (model.User, error) {
	var user model.User
	var date time.Time
	if err := row.Scan(&user.ID, &user.Name, &user.FirstName, &user.LastName,
		&date, &user.Picture, &user.Locale); err != nil {
		log.Println("Error reading user: ", err)
		return user, err
	}

	return user, nil
}

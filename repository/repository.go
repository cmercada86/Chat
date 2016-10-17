package repository

import (
	"Chat/model"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

const maxConnections = 20

var db *sql.DB

func NewRepository(user string, pass string, host string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=chat host=%s sslmode=disable ",
		user,
		pass,
		host)

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
		log.Println("Error adding or updating user: ", err)
	}

}

func AddChatMessage(room string, userID string, message string) {
	query := fmt.Sprintf(
		`
		INSERT INTO chat_table (uid, user_id, timestamp,room, message)
		VALUES (uuid_generate_v4(),'%s',NOW(),'%s','%s')

		;`, userID, room, message)

	if err := db.QueryRow(query).Scan(); err != nil && err != sql.ErrNoRows {
		log.Println("Error adding or updating user: ", err)
	}
	//log.Println("Room: ", room, " Message: ", message)
}

func GetChatMessages(room string) ([]model.Chat, error) {

	var messages []model.Chat
	query := fmt.Sprintf("SELECT * FROM chat_table WHERE room='%s';", room)

	rows, err := db.Query(query)
	if err != nil && err != sql.ErrNoRows {
		//log.Println("Error querying chat_table: ", err)
		return messages, err
	}

	for rows.Next() {
		var message model.Chat
		rows.Scan(&message)

		messages = append(messages, message)

	}

	return messages, nil

}

func GetRoomNames() ([]string, error) {
	var rooms []string

	rows, err := db.Query("SELECT DISCTINCT room from chat_table;")
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

func Close() {
	db.Close()
}

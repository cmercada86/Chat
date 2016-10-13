package repository

import (
	"Chat/model"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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

func Close() {
	db.Close()
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

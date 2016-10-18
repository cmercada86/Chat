package auth

import (
	"Chat/cache"
	"Chat/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	applicationName = "Grow Chat Auth"
)

type ClientInfo struct {
	ClientID string
	State    model.UUID
}

var config = &oauth2.Config{
	//ClientID:     clientID,
	//ClientSecret: clientSecret,
	// Scope determines which API calls you are authorized to make
	Scopes: []string{"https://www.googleapis.com/auth/plus.login",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
	// Use "postmessage" for the code-flow for server side apps
	RedirectURL: "postmessage",
}

func SetOAuth2Config(clientID string, clientSecret string) {
	config.ClientID = clientID
	config.ClientSecret = clientSecret
}

func ExchangeAuthCodeForUser(authCode string) (model.User, bool, error) {
	token, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		return model.User{}, false, err
	}

	client := config.Client(oauth2.NoContext, token)
	//userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return model.User{}, false, err
	}
	defer userInfo.Body.Close()
	data, _ := ioutil.ReadAll(userInfo.Body)

	googleID, _ := model.GetUserIDfromGoogleLogin(data)

	plusInfo, err := client.Get("https://www.googleapis.com/plus/v1/people/" + googleID)
	defer plusInfo.Body.Close()
	plusData, _ := ioutil.ReadAll(plusInfo.Body)

	return model.UserFromGooglePlusUser(plusData)

}

func CheckAuth(state model.UUID) (model.User, bool) {
	if !cache.Contains(state) {
		log.Println("invalid state!")
		return model.User{}, false
	}

	user := cache.Get(state)

	return user, !(user.ID == "")

}

// decodeIdToken takes an ID Token and decodes it to fetch the Google+ ID within
func decodeIDToken(idToken string) (gplusID string, err error) {

	var set string
	if idToken != "" {
		// Check that the padding is correct for a base64decode
		parts := strings.Split(idToken, ".")
		if len(parts) < 2 {
			return "", fmt.Errorf("Malformed ID token")
		}
		// Decode the ID token
		b, err := base64Decode(parts[1])
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
		err = json.Unmarshal(b, &set)
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
	}
	return set, nil
}
func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

package model

import "encoding/json"

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastName"`
	Picture   string `json:"picture"`
	Locale    string `json:"locale"`
}

type googleuser struct {
	FamilyName string `json:"family_name"`
	Gender     string `json:"gender"`
	GivenName  string `json:"given_name"`
	Locale     string `json:"locale"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	Profile    string `json:"profile"`
	Sub        string `json:"sub"`
}

func UserFromGoogleUser(info []byte) (User, error) {

	var gUser googleuser

	if err := json.Unmarshal(info, &gUser); err != nil {
		//log.Println(err)
		return User{}, err
	}

	return User{
		ID:        gUser.Sub,
		Name:      gUser.Name,
		FirstName: gUser.GivenName,
		LastName:  gUser.FamilyName,
		Picture:   gUser.Picture,
		Locale:    gUser.Locale,
	}, nil
}

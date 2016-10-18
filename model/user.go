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
	// FamilyName string `json:"family_name"`
	// Gender     string `json:"gender"`
	// GivenName  string `json:"given_name"`
	// Locale     string `json:"locale"`
	// Name       string `json:"name"`
	// Picture    string `json:"picture"`
	// Profile    string `json:"profile"`
	Sub string `json:"sub"`
}

type googleplus struct {
	Kind        string `json:"kind"`
	Etag        string `json:"etag"`
	Nickname    string `json:"nickname"`
	ObjectType  string `json:"objectType"`
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Name        struct {
		FamilyName string `json:"familyName"`
		GivenName  string `json:"givenName"`
	} `json:"name"`
	URL   string `json:"url"`
	Image struct {
		URL       string `json:"url"`
		IsDefault bool   `json:"isDefault"`
	} `json:"image"`
	IsPlusUser bool   `json:"isPlusUser"`
	Language   string `json:"language"`
	AgeRange   struct {
		Min int `json:"min"`
	} `json:"ageRange"`
	CircledByCount int  `json:"circledByCount"`
	Verified       bool `json:"verified"`
}

// func UserFromGoogleUser(info []byte) (User, error) {

// 	var gUser googleuser

// 	if err := json.Unmarshal(info, &gUser); err != nil {
// 		//log.Println(err)
// 		return User{}, err
// 	}

// 	return User{
// 		ID:        gUser.Sub,
// 		Name:      gUser.Name,
// 		FirstName: gUser.GivenName,
// 		LastName:  gUser.FamilyName,
// 		Picture:   gUser.Picture,
// 		Locale:    gUser.Locale,
// 	}, nil
// }
func GetUserIDfromGoogleLogin(info []byte) (string, error) {
	var gUser googleuser
	if err := json.Unmarshal(info, &gUser); err != nil {
		//log.Println(err)
		return "", err
	}
	return gUser.Sub, nil
}

func UserFromGooglePlusUser(info []byte) (User, bool, error) {
	var gPlus googleplus

	if err := json.Unmarshal(info, &gPlus); err != nil {
		//log.Println(err)
		return User{}, false, err
	}
	return User{
		ID:        gPlus.ID,
		Name:      gPlus.DisplayName,
		FirstName: gPlus.Name.GivenName,
		LastName:  gPlus.Name.FamilyName,
		Picture:   gPlus.Image.URL,
		Locale:    gPlus.Language,
	}, gPlus.IsPlusUser, nil
}

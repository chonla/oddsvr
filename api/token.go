package api

import (
	"encoding/json"
)

// TokenEx for doing token exchange with strava api
type TokenEx struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

// TokenRefresh for doing token refresh
type TokenRefresh struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

// Athlete is runner profile
type Athlete struct {
	ID                   uint32 `json:"id"`
	UserName             string `json:"username"`
	FirstName            string `json:"firstname"`
	LastName             string `json:"lastname"`
	City                 string `json:"city"`
	State                string `json:"state"`
	Country              string `json:"country"`
	Gender               string `json:"sex"`
	ProfilePicture       string `json:"profile"`
	ProfilePictureMedium string `json:"profile_medium"`
	Email                string `json:"email"`
}

// Token is token responded from strava
type Token struct {
	TokenType    string `json:"token_type"`
	Expiry       uint32 `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      `json:"athlete"`
}

// InvertedToken is token represented in database
type InvertedToken struct {
	ID uint32 `json:"_id"`
	*Token
}

func (t *Token) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

func (a *API) saveToken(t *InvertedToken) error {
	e := a.dbc.Replace("athlete", t.ID, t)
	return e
}

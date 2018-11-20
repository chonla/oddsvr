package api

import (
	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
)

// c represents customized of claims
type JWTClaims struct {
	ID          uint32 `json:"id"`
	StravaToken string `json:"strava_token"`
	jwt.StandardClaims
}

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

func (a *Athlete) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

func (a *API) saveToken(t *InvertedToken) error {
	e := a.dbc.Replace("athlete", t.ID, t)
	return e
}

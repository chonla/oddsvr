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
	Stats                `json:"stats"`
}

// Stats is running stats
type Stats struct {
	RecentRunTotals RunStats `json:"recent_run_totals"`
	AllRunTotals    RunStats `json:"all_run_totals"`
}

// RunStats is detailed of stats
type RunStats struct {
	Count         uint32  `json:"count"`
	Distance      float64 `json:"distance"`
	MovingTime    uint32  `json:"moving_time"`
	ElapsedTime   uint32  `json:"elapsed_time"`
	ElevationGain float64 `json:"elevation_gain"`
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

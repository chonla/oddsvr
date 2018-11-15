package api

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const (
	apiBase  = "https://www.strava.com/api/v3"
	apiOAuth = "https://www.strava.com/oauth/token"
)

// MeHandler returns information of myself
func (a *API) MeHandler(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	token := claims.StravaToken

	s := NewStrava()

	me, e := s.Me(token)
	if e != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	return c.String(http.StatusOK, fmt.Sprint(me))
}

// GatewayHandler handles request redirect_uri
func (a *API) GatewayHandler(c echo.Context) error {
	code := c.QueryParam("code")
	tokenEx := TokenEx{
		ClientID:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		Code:         code,
	}

	s := NewStrava()

	token, _ := s.TokenExchange(tokenEx)
	e := a.saveToken(&InvertedToken{
		ID:    token.ID,
		Token: token,
	})
	if e != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	claims := &JWTClaims{
		ID:          token.ID,
		StravaToken: token.AccessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, e := jwtToken.SignedString([]byte(a.config.JWTSecret))
	if e != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)

	c.SetCookie(cookie)

	return c.Redirect(http.StatusTemporaryRedirect, "http://localhost:4200/dashboard")
}

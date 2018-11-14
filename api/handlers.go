package api

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const (
	apiOAuth = "https://www.strava.com/oauth/token"
)

// GatewayHandler handles request redirect_uri
func (a *API) GatewayHandler(c echo.Context) error {
	code := c.QueryParam("code")
	tokenEx := TokenEx{
		ClientID:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		Code:         code,
	}

	token, _ := tokenExchange(tokenEx)
	e := a.saveToken(&InvertedToken{
		ID:    token.ID,
		Token: token,
	})
	if e != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["id"] = token.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, e := jwtToken.SignedString([]byte("OkCigarette"))
	if e != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)

	c.SetCookie(cookie)

	return c.Redirect(302, "http://localhost:4200/dashboard")
}

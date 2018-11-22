package api

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

const (
	apiBase  = "https://www.strava.com/api/v3"
	apiOAuth = "https://www.strava.com/oauth/token"
)

// VrJoinHandler engages user to a run
func (a *API) VrJoinHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	uid := claims.ID

	vr := NewVr()
	id := c.Param("id")
	if a.hasVr(id) {
		e := a.loadVr(id, vr)
		if e != nil {
			return c.JSON(http.StatusInternalServerError, e)
		}

		vr.Athletes = append(vr.Athletes, uid)

		a.saveVr(vr)
	} else {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, vr)
}

// VrGetHandler returns virtual run info
func (a *API) VrGetHandler(c echo.Context) error {
	vr := NewVr()
	id := c.Param("id")
	if a.hasVr(id) {
		e := a.loadVr(id, vr)
		if e != nil {
			return c.JSON(http.StatusInternalServerError, e)
		}
	} else {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, vr)
}

// VrCreationHandler creates a new virtual run
func (a *API) VrCreationHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	uid := claims.ID

	vr, e := NewVrFromContext(c)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
	}

	newID := bson.NewObjectId()
	vr.ID = newID
	vr.Athletes = []uint32{
		uid,
	}

	id, e := a.saveVr(vr)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}
	c.Response().Header().Add("Location", fmt.Sprintf("/vr/%s", id))

	return c.JSON(http.StatusCreated, vr)
}

// MeGetHandler returns information of myself
func (a *API) MeGetHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	token := claims.StravaToken

	s := NewStrava()

	me, e := s.Me(token)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}

	return c.JSON(http.StatusOK, me)
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
		return c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	claims := &JWTClaims{
		ID:          token.ID,
		StravaToken: token.AccessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, e := jwtToken.SignedString([]byte(a.config.JWTSecret))
	if e != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(3 * time.Hour)

	c.SetCookie(cookie)

	return c.Redirect(http.StatusTemporaryRedirect, "http://localhost:4200/dashboard")
}

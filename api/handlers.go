package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/kr/pretty"

	jwt "github.com/dgrijalva/jwt-go"
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

	dist, e := NewDistanceFromContext(c)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}

	vr := NewVr()
	id := c.Param("id")
	if a.hasVrByLink(id) {
		e := a.loadVrByLink(id, vr)
		if e != nil {
			return c.JSON(http.StatusInternalServerError, e)
		}

		for _, eng := range vr.Engagements {
			if eng.Athlete == uid {
				return c.JSON(http.StatusOK, vr)
			}
		}

		vr.Engagements = append(vr.Engagements,
			Engagement{
				Athlete:  uid,
				Distance: dist.Distance,
			})

		a.saveVr(vr)
	} else {
		return c.NoContent(http.StatusNotFound)
	}

	c.Response().Header().Add("Location", fmt.Sprintf("/vr/%s", vr.Link))
	c.Response().Header().Add("X-Join-Vr-ID", id)
	return c.JSON(http.StatusCreated, vr)
}

// VrGetAvailableHandler returns virtual run info
func (a *API) VrGetAvailableHandler(c echo.Context) error {
	vrs := []VirtualRun{}
	e := a.loadAvailableVr(&vrs)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}
	vrsum := []VirtualRunSummary{}
	for _, v := range vrs {
		vrsum = append(vrsum, VirtualRunSummary{
			ID:               v.ID,
			CreatedBy:        v.CreatedBy,
			CreatedDateTime:  v.CreatedDateTime,
			Period:           v.Period,
			Title:            v.Title,
			Link:             v.Link,
			EngagementsCount: len(v.Engagements),
		})
	}
	return c.JSON(http.StatusOK, vrsum)
}

// VrGetMineHandler returns virtual run info
func (a *API) VrGetMineHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	uid := claims.ID

	vrs := []VirtualRun{}
	e := a.loadMyVr(uid, &vrs)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}
	vrsum := []VirtualRunSummary{}
	for _, v := range vrs {
		vrsum = append(vrsum, VirtualRunSummary{
			ID:               v.ID,
			CreatedBy:        v.CreatedBy,
			CreatedDateTime:  v.CreatedDateTime,
			Period:           v.Period,
			Title:            v.Title,
			Link:             v.Link,
			EngagementsCount: len(v.Engagements),
		})
	}
	return c.JSON(http.StatusOK, vrsum)
}

// VrGetByLinkHandler returns virtual run info
func (a *API) VrGetByLinkHandler(c echo.Context) error {
	vr := NewVr()
	id := c.Param("id")
	if a.hasVrByLink(id) {
		e := a.loadVrByLink(id, vr)
		if e != nil {
			return c.JSON(http.StatusInternalServerError, e)
		}
	} else {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, vr)
}

// VrGetByPrivateLinkHandler returns virtual run info
func (a *API) VrGetByPrivateLinkHandler(c echo.Context) error {
	var uid uint32
	if c.Get("user") != nil {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JWTClaims)
		uid = claims.ID
		// pretty.Println("uid = %d", uid)
	}

	vr := NewVr()
	id := c.Param("id")
	if a.hasVrByLink(id) {
		e := a.loadVrByLink(id, vr)
		if e != nil {
			return c.JSON(http.StatusInternalServerError, e)
		}
	} else {
		return c.NoContent(http.StatusNotFound)
	}
	for _, eng := range vr.Engagements {
		pretty.Println(eng)
		pretty.Println(uid)
		if eng.Athlete == uid {
			break
		}
	}
	return c.JSON(http.StatusOK, vr)
}

// VrCreationHandler creates a new virtual run
func (a *API) VrCreationHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	uid := claims.ID

	vrc, e := NewVrFromContext(c)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}

	vr := NewVr()
	vr.ID = bson.NewObjectId()
	vr.Title = vrc.Title
	vr.Detail = vrc.Detail
	vr.Period = vrc.Period
	vr.Link = a.createSafeVrLink()
	vr.CreatedBy = uid
	vr.CreatedDateTime = time.Now().Format(time.RFC3339)

	id, e := a.saveVr(vr)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}
	c.Response().Header().Add("Location", fmt.Sprintf("/vr/%s", vr.Link))
	c.Response().Header().Add("X-New-Vr-ID", id)

	return c.JSON(http.StatusCreated, vr)
}

// MeGetHandler returns information of myself
func (a *API) MeGetHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	token := claims.StravaToken

	s := NewStrava(a.cache)

	me, e := s.Me(token)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}

	return c.JSON(http.StatusOK, me)
}

func (a *API) gateway(c echo.Context) error {
	code := c.QueryParam("code")
	tokenEx := TokenEx{
		ClientID:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		Code:         code,
	}

	s := NewStrava(a.cache)

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
	cookie.Path = "/"

	c.SetCookie(cookie)

	myIdCookie := new(http.Cookie)
	myIdCookie.Name = "me"
	myIdCookie.Value = fmt.Sprintf("%d", token.ID)
	myIdCookie.Expires = time.Now().Add(3 * time.Hour)
	myIdCookie.Path = "/"

	return nil
}

// GatewayAndGoToHandler handles request redirect_uri
func (a *API) GatewayAndGoToHandler(c echo.Context) error {
	id := c.Param("id")

	e := a.gateway(c)
	if e != nil {
		return e
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/vr/%s", a.config.FrontBase, id))
}

// GatewayHandler handles request redirect_uri
func (a *API) GatewayHandler(c echo.Context) error {
	e := a.gateway(c)
	if e != nil {
		return e
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/vr", a.config.FrontBase))
}

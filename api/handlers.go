package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kr/pretty"

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
	return c.JSON(http.StatusOK, vr)
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
			FromDate:         v.FromDate,
			ToDate:           v.ToDate,
			Name:             v.Name,
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
	vr.Joined = false
	return c.JSON(http.StatusOK, vr)
}

// VrGetByPrivateLinkHandler returns virtual run info
func (a *API) VrGetByPrivateLinkHandler(c echo.Context) error {
	var uid uint32
	if c.Get("user") != nil {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JWTClaims)
		uid = claims.ID
		pretty.Println("uid = %d", uid)
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
	vr.Joined = false
	for _, eng := range vr.Engagements {
		pretty.Println(eng)
		pretty.Println(uid)
		if eng.Athlete == uid {
			vr.Joined = true
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
	vr.Name = vrc.Name
	vr.FromDate = vrc.FromDate
	vr.ToDate = vrc.ToDate
	vr.Link = a.createSafeVrLink()

	newID := bson.NewObjectId()
	vr.ID = newID
	vr.Engagements = []Engagement{
		Engagement{
			Athlete:  uid,
			Distance: vrc.Distance,
		},
	}
	vr.CreatedBy = uid
	vr.Joined = true

	now := time.Now()
	vr.CreatedDateTime = now.Format(time.RFC3339)

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

func (a *API) gateway(c echo.Context) error {
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
	cookie.Path = "/"

	c.SetCookie(cookie)

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

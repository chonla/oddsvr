package api

import (
	"fmt"

	"github.com/chonla/oddsvr/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// API represents API type
type API struct {
	config *Config
	dbc    *db.Client
}

// Config represents configuration
type Config struct {
	DatabaseConnectionString string
	ClientID                 string
	ClientSecret             string
	JWTSecret                string
}

// NewAPI creates a new API
func NewAPI(conf *Config) (*API, error) {
	Info(fmt.Sprintf("connecting to database %s...", conf.DatabaseConnectionString))
	dbConnection, e := db.NewClient(conf.DatabaseConnectionString, "vr")
	if e != nil {
		return nil, e
	}
	return &API{
		config: conf,
		dbc:    dbConnection,
	}, nil
}

// Serve starts service
func (a *API) Serve(addr string) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.CORS())

	e.GET("/gateway", a.GatewayHandler)

	r := e.Group("/api")
	jwtConfig := middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte(a.config.JWTSecret),
	}
	r.Use(middleware.JWTWithConfig(jwtConfig))
	r.GET("/me", a.MeGetHandler)
	r.POST("/vr", a.VrCreationHandler)
	r.GET("/vr/:id", a.VrGetHandler)
	r.GET("/vr/:id/join", a.VrJoinHandler)

	Info(fmt.Sprintf("server is listening on %s", addr))
	e.Logger.Fatal(e.Start(addr))
}

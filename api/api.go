package api

import (
	"fmt"

	"github.com/chonla/oddsvr/db"
	"github.com/labstack/echo"
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

	e.GET("/gateway", a.GatewayHandler)

	Info(fmt.Sprintf("server is listening on %s", addr))
	e.Logger.Fatal(e.Start(addr))
}

package main

import (
	"fmt"
	"os"

	"github.com/chonla/oddsvr/api"
)

func main() {
	api.Info("odds virtual run api server")

	stravaClientID, e := env("ODDSVR_STRAVA_CLIENT_ID", "", "strava client id", true)
	if e != nil {
		api.Error(e.Error())
		os.Exit(1)
	}
	stravaClientSecret, e := env("ODDSVR_STRAVA_CLIENT_SECRET", "", "strava client secret", true)
	if e != nil {
		api.Error(e.Error())
		os.Exit(1)
	}
	jwtSecret, e := env("ODDSVR_JWT_SECRET", "", "jwt secret", true)
	if e != nil {
		api.Error(e.Error())
		os.Exit(1)
	}
	dbServer, _ := env("ODDSVR_DB", "127.0.0.1:27017", "database address", false)

	conf := &api.Config{
		DatabaseConnectionString: dbServer,
		ClientID:                 stravaClientID,
		ClientSecret:             stravaClientSecret,
		JWTSecret:                jwtSecret,
	}
	vr, e := api.NewAPI(conf)
	if e != nil {
		api.Error(fmt.Sprintf("unable to start service: %v", e))
	} else {
		boundAddress, _ := env("ODDSVR_ADDR", ":1323", "application address", false)
		vr.Serve(boundAddress)
	}
}

func env(key, defaultValue, name string, errorIfMissing bool) (string, error) {
	value, found := os.LookupEnv(key)
	if !found || value == "" {
		if errorIfMissing {
			return "", fmt.Errorf("seems like %s (%s) is missing from environment variables", name, key)
		}
		api.Info(fmt.Sprintf("%s is set to default value", name))
		api.Info(fmt.Sprintf("you can override this using %s environment variable", key))
		value = defaultValue
	}
	return value, nil
}

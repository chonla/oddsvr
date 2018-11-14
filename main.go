package main

import (
	"fmt"
	"os"

	"github.com/chonla/oddsvr/api"
)

func main() {
	api.Info("odds virtual run api server")

	dbServer, _ := env("ODDSVR_DB", "127.0.0.1:27017", "database address", false)

	conf := &api.Config{
		DatabaseConnectionString: dbServer,
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
	value := os.Getenv(key)
	if value == "" {
		if errorIfMissing {
			return "", fmt.Errorf("%s is missing from environment variables", key)
		}
		api.Info(fmt.Sprintf("%s is set to default value", name))
		api.Info(fmt.Sprintf("you can override this using %s environment variable", key))
		value = defaultValue
	}
	return value, nil
}

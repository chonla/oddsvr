package main

import (
	"fmt"
	"os"

	"github.com/chonla/oddsvr/api"
)

func main() {
	api.Info("odds virtual run api server")

	dbServer := os.Getenv("ODDSVR_DB")
	if dbServer == "" {
		api.Info("database is set to default address")
		api.Info("you can override default database address using ODDSVR_DB environment variable")
		dbServer = "127.0.0.1:27017"
	}

	conf := &api.Config{
		DatabaseConnectionString: dbServer,
	}
	vr, e := api.NewAPI(conf)
	if e != nil {
		api.Error(fmt.Sprintf("unable to start service: %v", e))
	} else {
		boundAddress := os.Getenv("ODDSVR_ADDR")
		if boundAddress == "" {
			api.Info("application is bound to default address")
			api.Info("you can override default listening address using ODDSVR_ADDR environment variable")
			boundAddress = ":1323"
		}

		vr.Serve(boundAddress)
	}
}

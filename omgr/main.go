package main

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/oapiclient"
	"apiscrape/omgr/data"
	"apiscrape/omgr/server"
)

func main() {
	logger := loggr.NewLogger()

	// store init
	logger.Log("initializing store")
	ac := oapiclient.NewOApiClient("oapi", 8080, logger, nil)
	st := data.NewStore(logger, ac)
	err := st.LoadConfigIntoStore()
	if err != nil {
		logger.Fatal("error on store init", err)
	}

	// oapi healthcheck
	logger.Log("testing oapi")
	tick, done := ac.WaitForHealth()
	for {
		complete := false
		select {
		case <-tick:
			logger.Log("checking oapi health..")
		case <-done:
			complete = true
			break
		}
		if complete {
			break
		}
	}
	logger.Log("oapi is online..")

	// get oapi data start time
	logger.Log("fetching start time from oapi")
	err = st.FetchOApiStartTime()
	if err != nil {
		logger.Fatal("unable to fetch start time", err)
	}
	logger.Log("start time is", st.ApiStartTime)

	// fetch locations concurrently
	logger.Log("fetching locations from oapi")
	err = st.FetchLocations()
	if err != nil {
		logger.Fatal("error while fetching locations", err)
	}

	// generating fetch jobs for workers
	logger.Log("generating fetch jobs for workers")
	st.GenerateFetchJobs()

	// starting server
	logger.Log("starting server")
	sv := server.NewApiServer(st, logger)
	err = sv.Start(":8080")
	if err != nil {
		logger.Fatal(err)
	}
}

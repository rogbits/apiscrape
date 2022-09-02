package main

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/oapiclient"
	"apiscrape/lib/omgrclient"
)

import (
	"apiscrape/owrk/data"
)

func main() {
	logger := loggr.NewLogger()

	// store init
	logger.Log("initializing store")
	mc := omgrclient.NewOMgrClient("omgr", 8080, logger, nil)
	ac := oapiclient.NewOApiClient("oapi", 8080, logger, nil)
	st := data.NewOWrkStore(logger, mc, ac)
	err := st.LoadConfigIntoStore()
	if err != nil {
		logger.Fatal("error on store init", err)
	}

	// omgr healthcheck
	logger.Log("checking omgr api")
	tick, done := mc.WaitForHealth()
	for {
		complete := false
		select {
		case <-tick:
			logger.Log("checking omgr api health..")
		case <-done:
			complete = true
			break
		}
		if complete {
			break
		}
	}
	logger.Log("omgr api is online..")

	// oapi healthcheck
	logger.Log("checking oapi")
	tick, done = mc.WaitForHealth()
	for {
		complete := false
		select {
		case <-tick:
			logger.Log("checking  oapi health..")
		case <-done:
			complete = true
			break
		}
		if complete {
			break
		}
	}
	logger.Log("oapi is online..")

	st.DequeueJobs()
}

package main

import (
	"apiscrape/lib/loggr"
	"apiscrape/oapi/data"
	"apiscrape/oapi/server"
)

func main() {
	logger := loggr.NewLogger()

	// store init
	logger.Log("initializing store")
	st := data.NewOApiStore(logger)
	err := st.Init()
	if err != nil {
		logger.Fatal("error on store init", err)
	}

	// view of generated data
	st.PrintLocationTicketTotals()
	st.PrintLocationTopItems()
	logger.Log("store start time:", st.StartTime)
	logger.Log("store end time:", st.StartTime+30*24*60*60)

	// server
	sv := server.NewApiServer(st, logger)
	err = sv.Start(":8080")
	if err != nil {
		logger.Fatal(err)
	}
}

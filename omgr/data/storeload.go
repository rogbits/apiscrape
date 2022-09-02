package data

import (
	"apiscrape/lib/oapiclient"
	"bufio"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"
)

func (store *Store) LoadConfigIntoStore() error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.New("could not get working directory")
	}

	fp := path.Join(wd, "config")
	file, err := os.Open(fp)
	if err != nil {
		return errors.New("could not open config file")
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		switch {
		case strings.HasPrefix(word, "oapi_batch_size="):
			value := strings.TrimLeft(word, "oapi_batch_size=")
			store.OApiBatchSize, _ = strconv.ParseInt(value, 10, 64)
		case strings.HasPrefix(word, "oapi_concurrent_requests="):
			value := strings.TrimLeft(word, "oapi_concurrent_requests=")
			store.OApiConcurrentRequests, _ = strconv.ParseInt(value, 10, 64)
		case strings.HasPrefix(word, "job_window_in_minutes="):
			value := strings.TrimLeft(word, "job_window_in_minutes=")
			store.JobTimeWindowInMinutes, _ = strconv.ParseInt(value, 10, 64)
		}
	}

	err = scanner.Err()
	return err
}

func (store *Store) GenerateFetchJobs() {
	i := int64(0)
	for key := range store.LocationsById {
		loc := store.LocationsById[key]
		start := store.ApiStartTime
		end := start + 30*24*60*60
		interval := 60 * store.JobTimeWindowInMinutes
		for j := start; j < end; j += interval {
			i++
			windowStart := j
			windowEnd := j + interval - 1
			if j+interval == end {
				// last job for location
				windowEnd = j + interval
			}
			fj := oapiclient.NewTicketFetchJob(i, loc.Id, loc.Name, windowStart, windowEnd)
			store.TicketFetchJobs = append(store.TicketFetchJobs, fj)
		}
	}
	store.Logger.Log(
		"total jobs generated", i,
		"for locations totalling", len(store.LocationsById),
	)
}

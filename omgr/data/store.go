package data

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/models"
	"apiscrape/lib/oapiclient"
	"sync"
)

type Store struct {
	OApiBatchSize      int64
	LocationsById      map[string]*models.Location
	LocationItemCounts map[string]map[string]uint64
	LocationTopItems   map[string][]string

	ApiStartTime           int64
	JobTimeWindowInMinutes int64
	TicketFetchJobs        []*oapiclient.FetchTicketJob
	JobsInFlight           map[int64]*oapiclient.FetchTicketJob

	OApiClient             *oapiclient.OApiClient
	OApiConcurrentRequests int64

	Mu     sync.Mutex
	Logger *loggr.Logger
}

func NewStore(logger *loggr.Logger, oApiClient *oapiclient.OApiClient) *Store {
	st := new(Store)
	st.LocationsById = map[string]*models.Location{}
	st.LocationItemCounts = map[string]map[string]uint64{}
	st.LocationTopItems = map[string][]string{}
	st.JobsInFlight = map[int64]*oapiclient.FetchTicketJob{}
	st.Logger = logger
	st.OApiClient = oApiClient
	return st
}

func (store *Store) FetchOApiStartTime() error {
	startTime, err := store.OApiClient.FetchStartTime()
	if err != nil {
		return err
	}
	store.ApiStartTime = startTime
	return nil
}

func (store *Store) FetchLocations() error {
	pageStart := int64(1)
	pageLimit := store.OApiBatchSize
	errors := make(chan error, store.OApiConcurrentRequests)
	for {
		hasMore := false
		var wg sync.WaitGroup
		wg.Add(int(store.OApiConcurrentRequests))
		for i := int64(0); i < store.OApiConcurrentRequests; i++ {
			ii := i
			go func() {
				// make request
				lq := models.NewBlankLocationQuery()
				lq.PageStart = pageLimit*ii + pageStart
				lq.PageLimit = pageLimit
				lq, err := store.OApiClient.FetchLocations(lq)
				if err != nil {
					errors <- err
					wg.Done()
					return
				}
				// populate locations from resp
				store.Mu.Lock()
				for _, loc := range lq.Result.Locations {
					store.LocationsById[loc.Id] = loc
				}
				store.Mu.Unlock()
				// trigger next batch
				// if next link on last page
				if ii == store.OApiConcurrentRequests-1 &&
					lq.HalResp.Links.Next.Href != "" {
					nextPageStart, err := lq.HalResp.GetNextPageStart()
					if err != nil {
						errors <- err
						wg.Done()
						return
					}
					if nextPageStart > pageStart {
						pageStart = nextPageStart
					}
					hasMore = true
				}
				wg.Done()
			}()
		}
		wg.Wait()
		if len(errors) > 0 {
			return <-errors
		}
		if !hasMore {
			break
		}
	}
	return nil
}

func (store *Store) GetTopItemsByLocationId(locId string) ([]string, uint64) {
	topItems := store.LocationTopItems[locId]
	return topItems, store.LocationItemCounts[locId][topItems[0]]
}

func (store *Store) ProcessJobResult(fj *oapiclient.FetchTicketJob) {
	_, exists := store.LocationItemCounts[fj.LocationId]
	if !exists {
		store.LocationItemCounts[fj.LocationId] = map[string]uint64{}
	}
	locId := fj.LocationId
	for _, item := range fj.Items {
		currentItemCount := store.LocationItemCounts[locId][item]
		store.Logger.Log("location", locId, "adding", item, "with latest total", currentItemCount+1)
		if len(store.LocationTopItems[locId]) == 0 {
			store.LocationTopItems[locId] = []string{item}
			store.LocationItemCounts[locId][item] += 1
			continue
		}

		aTopItemInLocation := store.LocationTopItems[locId][0]
		currentTopItemCountForLocation := store.LocationItemCounts[locId][aTopItemInLocation]
		switch {
		case currentItemCount+1 == currentTopItemCountForLocation:
			store.LocationItemCounts[locId][item] += 1
			store.LocationTopItems[locId] = append(store.LocationTopItems[locId], item)
		case currentItemCount+1 > currentTopItemCountForLocation:
			store.LocationItemCounts[locId][item] += 1
			store.LocationTopItems[locId] = []string{item}
		default:
			store.LocationItemCounts[locId][item] += 1
		}
	}
}

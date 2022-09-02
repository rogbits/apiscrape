package data

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/models"
	"apiscrape/lib/tools"
	"fmt"
	"sort"
	"strings"
	"time"
)

type OApiStore struct {
	NumLocations          int64
	MinTicketsPerLocation int64
	MaxTicketsPerLocation int64
	MinItemsPerTicket     int64
	MaxItemsPerTicket     int64

	StartTime int64

	LocationList        models.LocationList
	LocationListPos     map[string]int64
	LocationTickets     map[string]*models.TicketList
	LocationTicketTimes map[string][]int64

	ItemList           models.ItemList
	LocationItemCount  map[string]int64
	LocationItemTotals map[string]map[string]int64

	Logger *loggr.Logger
}

func NewOApiStore(logger *loggr.Logger) *OApiStore {
	s := new(OApiStore)
	s.LocationListPos = map[string]int64{}
	s.LocationTickets = map[string]*models.TicketList{}
	s.LocationTicketTimes = map[string][]int64{}
	s.LocationItemTotals = map[string]map[string]int64{}
	s.LocationItemCount = map[string]int64{}
	s.Logger = logger
	return s
}

func (store *OApiStore) GetLocationById(id string) *models.Location {
	index, exists := store.LocationListPos[id]
	if !exists {
		return nil
	}
	return store.LocationList.GetLocation(index)
}

func (store *OApiStore) GetItemCountForLocation(loc *models.Location, item *models.Item) int64 {
	return store.LocationItemTotals[loc.Id][item.Name]
}

func (store *OApiStore) GetTopItemsForLocation(loc *models.Location) ([]*models.Item, int64) {
	var topItems []*models.Item
	var topCount int64
	for _, item := range store.ItemList.Items {
		count := store.GetItemCountForLocation(loc, item)
		switch {
		case count > topCount:
			topItems = []*models.Item{}
			topItems = append(topItems, item)
			topCount = count
		case count == topCount:
			topItems = append(topItems, item)
		}
	}
	return topItems, topCount
}

func (store *OApiStore) GetLocationsByLocationQuery(lq *models.LocationQuery) *models.LocationQuery {
	i := lq.PageStart - 1
	j := i + lq.PageLimit
	l := int64(len(store.LocationList.Locations))
	if j > l {
		j = l
	}

	locationList := models.NewLocationList()
	if i < l {
		locationList.Locations = store.LocationList.Locations[i:j]
	}
	lq.Result = locationList
	lq.PagedCount = int64(len(locationList.Locations))
	lq.TotalNumRecords = int64(len(store.LocationList.Locations))
	return lq
}

func (store *OApiStore) GetLocationTicketByTicketQuery(tq *models.TicketQuery) *models.TicketQuery {
	location := store.GetLocationById(tq.Location.Id)
	locationTickets := store.LocationTickets[location.Id]
	locationTicketTimes := store.LocationTicketTimes[location.Id]

	rs := tools.NewRangeSearch(locationTicketTimes, tq.CloseStart, tq.CloseEnd)
	rs.Execute()
	left := rs.Left
	right := rs.Right
	var startIndex, endIndex int64
	if right == -1 {
		endIndex = -1
	} else {
		startIndex = left + tq.PageStart - 1
		endIndex = startIndex + tq.PageLimit
		if endIndex-right > 1 {
			endIndex = right + 1
		}
	}
	ticketList := models.NewTicketList()
	if endIndex > 0 {
		ticketList.Tickets = locationTickets.Tickets[startIndex:endIndex]
	}
	for _, ticket := range ticketList.Tickets {
		store.Logger.Log(
			tq,
			"including ticket", ticket.Id,
			"totalInQuery", right-left+1,
		)
		//store.Logger.Log(locationTicketTimes)
	}
	tq.Result = ticketList
	tq.PagedCount = int64(len(ticketList.Tickets))
	tq.TotalNumRecords = right - left + 1
	if tq.TotalNumRecords < 0 {
		tq.TotalNumRecords = 0
	}
	return tq
}

func (store *OApiStore) PrintConfiguration() {
	store.Logger.Log("num locations:", store.NumLocations)
	store.Logger.Log("ticket range per location:", fmt.Sprintf("%d-%d",
		store.MinTicketsPerLocation, store.MaxTicketsPerLocation))
	store.Logger.Log("item range per ticket:", fmt.Sprintf("%d-%d",
		store.MinItemsPerTicket, store.MaxItemsPerTicket))
	store.Logger.Log("num items added:", len(store.ItemList.Items))
}

func (store *OApiStore) PrintLocationTicketTotals() {
	for key := range store.LocationTickets {
		TicketList := store.LocationTickets[key]
		loc := store.GetLocationById(key)
		store.Logger.Log("total tickets for", loc.Name, loc.Id, ":",
			len(TicketList.Tickets))
	}
}

func (store *OApiStore) PrintLocationTopItems() {
	for _, loc := range store.LocationList.Locations {
		items, count := store.GetTopItemsForLocation(loc)
		var s []string
		for _, item := range items {
			s = append(s, item.Name)
		}
		store.Logger.Log("top items for", loc.Name, loc.Id, strings.Join(s, " "), "with", count)
	}
}

func (store *OApiStore) PrintLocationItemCounts(alphabetical bool) {
	for _, loc := range store.LocationList.Locations {
		var names []string
		var counts []int64
		for _, item := range store.ItemList.Items {
			count := store.GetItemCountForLocation(loc, item)
			names = append(names, item.Name)
			counts = append(counts, count)
		}
		store.Logger.Log("item total for", loc.Name, loc.Id)
		if alphabetical {
			sort.Sort(tools.NewTwoSliceSort(names, counts))
			for i, _ := range names {
				store.Logger.Log(loc.Name, loc.Id, names[i], counts[i])
			}
		} else {
			sort.Sort(tools.NewTwoSliceSort(counts, names))
			for i, _ := range names {
				store.Logger.Log(loc.Name, loc.Id, counts[i], names[i])
			}
		}
	}
}

func (store *OApiStore) PrintTicketsPerDay() {
	for _, loc := range store.LocationList.Locations {
		tickets := store.LocationTickets[loc.Id].Tickets
		dayKeys := make([]string, 0, 30)
		dayCount := map[string]int{}
		for _, ticket := range tickets {
			t := time.Unix(ticket.ClosedAt, 0)
			t = tools.RoundDownToDay(t)
			_, exists := dayCount[t.String()]
			if !exists {
				dayKeys = append(dayKeys, t.String())
				dayCount[t.String()] = 0
			}
			dayCount[t.String()]++
		}
		sort.Strings(dayKeys)
		for _, key := range dayKeys {
			store.Logger.Log("location",
				loc.Name, loc.Id, key, dayCount[key])
		}
	}
}

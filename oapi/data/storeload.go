package data

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

import (
	"apiscrape/lib/models"
	"apiscrape/lib/tools"
)

func (store *OApiStore) Init() error {
	err := store.LoadConfigIntoStore()
	if err != nil {
		return err
	}
	store.GenerateLocations()
	store.GenerateTicketsForLocation()
	return nil
}

func (store *OApiStore) LoadConfigIntoStore() error {
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
	itemId := int64(0)

	for scanner.Scan() {
		word := scanner.Text()
		switch {
		case strings.HasPrefix(word, "num_locations="):
			value := strings.TrimLeft(word, "num_locations=")
			store.NumLocations, _ = strconv.ParseInt(value, 10, 64)
		case strings.HasPrefix(word, "min_tickets_per_location="):
			value := strings.TrimLeft(word, "min_tickets_per_location=")
			store.MinTicketsPerLocation, _ = strconv.ParseInt(value, 10, 64)
		case strings.HasPrefix(word, "max_tickets_per_location="):
			value := strings.TrimLeft(word, "max_tickets_per_location=")
			store.MaxTicketsPerLocation, _ = strconv.ParseInt(value, 10, 64)
		case strings.HasPrefix(word, "min_items_per_ticket="):
			value := strings.TrimLeft(word, "min_items_per_ticket=")
			store.MinItemsPerTicket, _ = strconv.ParseInt(value, 10, 64)
		case strings.HasPrefix(word, "max_items_per_ticket="):
			value := strings.TrimLeft(word, "max_items_per_ticket=")
			store.MaxItemsPerTicket, _ = strconv.ParseInt(value, 10, 64)
		case strings.HasPrefix(word, "items:"):
			continue
		default:
			item := models.NewItem()
			item.Name = word
			itemId++
			item.Id = itemId
			store.ItemList.AddItem(item)
		}
	}

	err = scanner.Err()
	return err
}

func (store *OApiStore) GenerateLocations() {
	for i := int64(0); i < store.NumLocations; i++ {
		locName := fmt.Sprintf("location%d", 10000+i)
		loc := models.NewLocation(locName)
		store.LocationList.AddLocation(loc)
		store.LocationTicketTimes[loc.Id] = []int64{}
		store.LocationTickets[loc.Id] = models.NewTicketList()
		store.LocationItemTotals[loc.Id] = map[string]int64{}
		store.LocationItemCount[loc.Id] = 0
		store.LocationListPos[loc.Id] = i
	}
}

func (store *OApiStore) GenerateTicketsForLocation() {
	now := tools.GetNowToTheMinute()
	thirtyDaysAgo := now.Add(-30 * 24 * time.Hour)
	store.StartTime = thirtyDaysAgo.Unix()

	for _, loc := range store.LocationList.Locations {
		numTickets := tools.GetRandomInt64(
			store.MinTicketsPerLocation,
			store.MaxTicketsPerLocation,
		)
		for i := int64(1); i <= numTickets; i++ {
			ticketTime := tools.GetRandomInt64(
				thirtyDaysAgo.Unix(),
				now.Unix(),
			)
			store.LocationTicketTimes[loc.Id] = tools.BinInsertInt64(
				store.LocationTicketTimes[loc.Id],
				ticketTime,
			)
		}
		for i, epoch := range store.LocationTicketTimes[loc.Id] {
			ticket := store.GenerateTicket(loc, int64(i+1), epoch)
			store.LocationTickets[loc.Id].AddTicket(ticket)
		}
	}
}

func (store *OApiStore) GenerateTicket(loc *models.Location, id, epoch int64) *models.Ticket {
	ticket := models.NewTicket()
	ticket.Id = id
	ticket.TicketNumber = ticket.Id
	ticket.Open = false
	ticket.Void = false
	ticket.OpenedAt = epoch
	ticket.ClosedAt = epoch
	numItems := tools.GetRandomInt64(
		store.MinItemsPerTicket,
		store.MaxItemsPerTicket,
	)

	for i := int64(1); i <= numItems; i++ {
		index := tools.GetRandomInt(0, len(store.ItemList.Items))
		item := store.ItemList.GetItem(int64(index))
		itemCopy := models.CopyItem(item)
		itemCopy.Quantity = 1
		ticket.Embedded.Items = append(ticket.Embedded.Items, itemCopy)

		store.LocationItemCount[loc.Id] += 1
		store.LocationItemTotals[loc.Id][item.Name] += 1
	}

	return ticket
}

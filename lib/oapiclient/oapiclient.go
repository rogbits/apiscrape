package oapiclient

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/models"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type OApiClient struct {
	Host       string
	Port       int
	Logger     *loggr.Logger
	HttpClient *http.Client
}

func NewOApiClient(host string, port int, logger *loggr.Logger, client *http.Client) *OApiClient {
	c := new(OApiClient)
	c.Host = host
	c.Port = port
	c.Logger = logger
	c.HttpClient = client
	if c.HttpClient == nil {
		c.HttpClient = &http.Client{}
	}
	return c
}

func (ac *OApiClient) WaitForHealth() (tick chan bool, done chan bool) {
	tick = make(chan bool)
	done = make(chan bool, 1)
	go func() {
		for {
			tick <- true
			err := ac.GetHealth()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		done <- true
	}()
	return tick, done
}

func (ac *OApiClient) GetHealth() error {
	url := fmt.Sprintf("http://%s:%d/1.0/health", ac.Host, ac.Port)
	resp, err := ac.HttpClient.Get(url)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	resText := scanner.Text()
	if resText == "ok" {
		return nil
	}
	return errors.New("oapi failed healthcheck")
}

func (ac *OApiClient) FetchStartTime() (int64, error) {
	url := fmt.Sprintf("http://%s:%d/1.0/start", ac.Host, ac.Port)
	resp, err := ac.HttpClient.Get(url)
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	startTime := scanner.Text()
	converted, err := strconv.ParseInt(startTime, 10, 64)
	if err != nil {
		return 0, err
	}
	return converted, nil
}

func (ac *OApiClient) FetchLocations(lq *models.LocationQuery) (*models.LocationQuery, error) {
	url := fmt.Sprintf(
		"http://%s:%d/1.0/locations?start=%d&limit=%d",
		ac.Host,
		ac.Port,
		lq.PageStart,
		lq.PageLimit,
	)
	resp, err := ac.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	halResp := models.NewHalResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(halResp)
	if err != nil {
		return nil, err
	}
	locationList, err := halResp.DecodeLocationsList()
	if err != nil {
		return nil, err
	}
	lq.Result = locationList
	lq.HalResp = halResp
	return lq, nil
}

func (ac *OApiClient) FetchTickets(tq *models.TicketQuery) (*models.TicketQuery, error) {
	url := fmt.Sprintf(
		"http://%s:%d/1.0/locations/%s/tickets?start=%d&limit=%d&where=%s&jobId=%d",
		ac.Host,
		ac.Port,
		tq.Location.Id,
		tq.PageStart,
		tq.PageLimit,
		tq.GenerateWhereParam(),
		tq.JobId,
	)
	resp, err := ac.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	halResp := models.NewHalResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(halResp)
	if err != nil {
		return nil, err
	}
	ticketsList, err := halResp.DecodeTicketsList()
	if err != nil {
		return nil, err
	}
	tq.Result = ticketsList
	tq.PagedCount = halResp.Count
	tq.HalResp = halResp
	return tq, nil
}

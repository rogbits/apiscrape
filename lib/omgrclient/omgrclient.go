package omgrclient

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/oapiclient"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type OMgrClient struct {
	Host string
	Port int

	Logger     *loggr.Logger
	HttpClient *http.Client
}

func NewOMgrClient(host string, port int, logger *loggr.Logger, client *http.Client) *OMgrClient {
	c := new(OMgrClient)
	c.Host = host
	c.Port = port
	c.Logger = logger
	c.HttpClient = client
	if c.HttpClient == nil {
		c.HttpClient = &http.Client{}
	}
	return c
}

func (mc *OMgrClient) GetHealth() error {
	url := fmt.Sprintf("http://%s:%d/1.0/health", mc.Host, mc.Port)
	resp, err := mc.HttpClient.Get(url)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	resText := scanner.Text()
	if resText == "ok" {
		return nil
	}
	return errors.New("omgr failed healthcheck")
}

func (mc *OMgrClient) WaitForHealth() (chan bool, chan bool) {
	tick := make(chan bool)
	done := make(chan bool, 1)
	go func() {
		for {
			tick <- true
			err := mc.GetHealth()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		done <- true
	}()
	return tick, done
}

func (mc *OMgrClient) GetJob() (*oapiclient.FetchTicketJob, error) {
	url := fmt.Sprintf("http://%s:%d/1.0/_dequeue", mc.Host, mc.Port)
	resp, err := mc.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	fj := oapiclient.NewTicketFetchJob(0, "", "", 0, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(fj)
	if err != nil {
		return nil, err
	}
	return fj, nil
}

func (mc *OMgrClient) UpdateJob(fj *oapiclient.FetchTicketJob) error {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(fj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s:%d/1.0/jobs", mc.Host, mc.Port)
	resp, err := mc.HttpClient.Post(url, "application/json", b)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("unexpected status code")
	}
	return nil
}

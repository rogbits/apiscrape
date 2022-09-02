package models

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type HalResp struct {
	Embedded interface{} `json:"_embedded,inline"`
	Links    Links       `json:"_links"`
	Count    int64       `json:"count"`
	Limit    int64       `json:"limit"`
}

type Links struct {
	Self Link
	Next Link
	Prev Link
}

type Link struct {
	Href string
	Type string
}

func NewHalResp() *HalResp {
	return new(HalResp)
}

func (hr *HalResp) DecodeLocationsList() (*LocationList, error) {
	m := hr.Embedded.(map[string]interface{})
	marshalled, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	locationList := NewLocationList()
	err = json.Unmarshal(marshalled, locationList)
	if err != nil {
		return nil, err
	}
	return locationList, nil
}

func (hr *HalResp) DecodeTicketsList() (*TicketList, error) {
	m := hr.Embedded.(map[string]interface{})
	marshalled, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	ticketList := NewTicketList()
	err = json.Unmarshal(marshalled, ticketList)
	if err != nil {
		return nil, err
	}
	return ticketList, nil
}

func (hr *HalResp) GetNextPageStart() (int64, error) {
	nextUrl, err := url.ParseRequestURI(hr.Links.Next.Href)
	if err != nil {
		return 0, err
	}
	start := nextUrl.Query().Get("start")
	if start == "" {
		return 0, errors.New("missing start param")
	}
	converted, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		return 0, err
	}
	return converted, nil
}

func (hr *HalResp) GetNextPageLimit() (int64, error) {
	nextUrl, err := url.ParseRequestURI(hr.Links.Next.Href)
	if err != nil {
		return 0, err
	}
	start := nextUrl.Query().Get("limit")
	if start == "" {
		return 0, errors.New("missing limit param")
	}
	converted, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		return 0, err
	}
	return converted, nil
}

func (hr *HalResp) GetPrevPageStart() (int64, error) {
	nextUrl, err := url.ParseRequestURI(hr.Links.Prev.Href)
	if err != nil {
		return 0, err
	}
	start := nextUrl.Query().Get("start")
	if start == "" {
		return 0, errors.New("missing start param")
	}
	converted, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		return 0, err
	}
	return converted, nil
}

func (hr *HalResp) GetPrevPageLimit() (int64, error) {
	nextUrl, err := url.ParseRequestURI(hr.Links.Prev.Href)
	if err != nil {
		return 0, err
	}
	start := nextUrl.Query().Get("limit")
	if start == "" {
		return 0, errors.New("missing limit param")
	}
	converted, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		return 0, err
	}
	return converted, nil
}

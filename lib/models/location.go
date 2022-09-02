package models

import (
	"apiscrape/lib/tools"
)

type Location struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewLocation(name string) *Location {
	l := new(Location)
	l.Id = tools.GenerateId(8)
	l.Name = name
	return l
}

type LocationList struct {
	Locations []*Location `json:"locations"`
}

func NewLocationList() *LocationList {
	ll := new(LocationList)
	ll.Locations = []*Location{}
	return ll
}

func (ls *LocationList) AddLocation(loc *Location) {
	ls.Locations = append(ls.Locations, loc)
}

func (ls *LocationList) GetLocation(i int64) *Location {
	return ls.Locations[i]
}

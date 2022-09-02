package models

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PaginatedQuery interface {
	GetPageStart() int64
	GetPageLimit() int64
	GetPagedCount() int64
	GetTotalNumRecords() int64
	GetResourceType() string
}

type TicketQuery struct {
	Location   *Location
	IsOpen     bool
	CloseStart int64
	CloseEnd   int64

	PageStart  int64
	PageLimit  int64
	PagedCount int64

	Result          *TicketList
	TotalNumRecords int64
	JobId           int64

	HalResp *HalResp
}

func NewBlankTicketQuery() *TicketQuery {
	return new(TicketQuery)
}

func NewTicketQuery(q url.Values) (*TicketQuery, error) {
	tq := new(TicketQuery)
	where := q.Get("where")
	pageStart := q.Get("start")
	pageLimit := q.Get("limit")

	// not a legit parse
	r1, err := regexp.Compile(`open,(true|false)`)
	if err != nil {
		return nil, err
	}
	r2, err := regexp.Compile(`gte\(closed_at,(\d+)\)`)
	if err != nil {
		return nil, err
	}
	r3, err := regexp.Compile(`lte\(closed_at,(\d+)\)`)
	if err != nil {
		return nil, err
	}

	open := r1.FindStringSubmatch(where)
	if len(open) > 1 {
		tq.IsOpen = strings.ToLower(open[1]) == "true"
	}
	closeStart := r2.FindStringSubmatch(where)
	if len(closeStart) > 1 {
		tq.CloseStart, err = strconv.ParseInt(closeStart[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}
	closeEnd := r3.FindStringSubmatch(where)
	if len(closeEnd) > 1 {
		tq.CloseEnd, err = strconv.ParseInt(closeEnd[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	tq.PageStart, _ = strconv.ParseInt(pageStart, 10, 64)
	if tq.PageStart == 0 {
		tq.PageStart = 1
	}
	tq.PageLimit, _ = strconv.ParseInt(pageLimit, 10, 64)
	if tq.PageLimit == 0 {
		tq.PageStart = 1
	}
	return tq, nil
}

func (tq *TicketQuery) String() string {
	locationKey := "nil"
	if tq.Location != nil {
		locationKey = tq.Location.Name + " " + tq.Location.Id
	}
	return fmt.Sprintf(
		"query: location:%s, open:%t, from:%s, to:%s, startAt: %d, limit: %d",
		locationKey,
		tq.IsOpen,
		time.Unix(int64(tq.CloseStart), 0).Format("2006-01-02T15:04:05"),
		time.Unix(int64(tq.CloseEnd), 0).Format("2006-01-02T15:04:05"),
		tq.PageStart,
		tq.PageLimit,
	)
}

func (tq *TicketQuery) GetPageStart() int64 {
	return tq.PageStart
}

func (tq *TicketQuery) GetPageLimit() int64 {
	return tq.PageLimit
}

func (tq *TicketQuery) GetPagedCount() int64 {
	return tq.PagedCount
}

func (tq *TicketQuery) GetTotalNumRecords() int64 {
	return tq.TotalNumRecords
}

func (tq *TicketQuery) GetResourceType() string {
	return "ticket"
}

func (tq *TicketQuery) GenerateWhereParam() string {
	return fmt.Sprintf(
		"and(eq(open,false),gte(closed_at,%d),lte(closed_at,%d))",
		tq.CloseStart,
		tq.CloseEnd,
	)
}

type LocationQuery struct {
	PageStart  int64
	PageLimit  int64
	PagedCount int64

	Result          *LocationList
	TotalNumRecords int64
	HalResp         *HalResp
}

func NewBlankLocationQuery() *LocationQuery {
	return new(LocationQuery)
}

func NewLocationQuery(q url.Values) (*LocationQuery, error) {
	lq := new(LocationQuery)
	pageStart := q.Get("start")
	pageLimit := q.Get("limit")
	lq.PageStart, _ = strconv.ParseInt(pageStart, 10, 64)
	if lq.PageStart == 0 {
		lq.PageStart = 1
	}
	lq.PageLimit, _ = strconv.ParseInt(pageLimit, 10, 64)
	if lq.PageLimit == 0 {
		lq.PageLimit = 1
	}

	return lq, nil
}

func (lq *LocationQuery) String() string {
	return fmt.Sprintf(
		"query:locations, start:%s, limit:%s",
		lq.PageStart,
		lq.PageLimit,
	)
}

func (lq *LocationQuery) GetPageStart() int64 {
	return lq.PageStart
}

func (lq *LocationQuery) GetPageLimit() int64 {
	return lq.PageLimit
}

func (lq *LocationQuery) GetPagedCount() int64 {
	return lq.PagedCount
}

func (lq *LocationQuery) GetTotalNumRecords() int64 {
	return lq.TotalNumRecords
}

func (lq *LocationQuery) GetResourceType() string {
	return "location"
}

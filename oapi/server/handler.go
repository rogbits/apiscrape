package server

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/models"
	"apiscrape/oapi/data"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type Handler struct {
	Logger *loggr.Logger
	Store  *data.OApiStore

	mu sync.Mutex
}

func NewHandler(store *data.OApiStore, logger *loggr.Logger) *Handler {
	h := new(Handler)
	h.Logger = logger
	h.Store = store
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Log("incoming", r.RequestURI)

	u, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		h.Logger.Log("error: bad URI", err)
		http.Error(w, "bad URI", http.StatusBadRequest)
		return
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		h.Logger.Log("error: bad query string", err)
		http.Error(w, "bad query string", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(u.Path, "/1.0") {
		http.Error(w, "missing version", http.StatusBadRequest)
		return
	}

	path := strings.TrimPrefix(u.Path, "/1.0")

	switch {
	case strings.HasPrefix(path, "/start"):
		startTime := strconv.FormatInt(h.Store.StartTime, 10)
		_, err := w.Write([]byte(startTime))
		if err != nil {
			h.Logger.Log("error: writing start time")
			http.Error(w, "error writing start time", http.StatusInternalServerError)
		}
		return
	case strings.HasPrefix(path, "/health"):
		_, err := w.Write([]byte("ok"))
		if err != nil {
			h.Logger.Log("error: writing health")
			http.Error(w, "error writing health", http.StatusInternalServerError)
		}
		return
	case strings.HasPrefix(path, "/locations/"):
		path = strings.TrimPrefix(path, "/locations/")
		if len(path) < 8 || path[8] != '/' {
			h.Logger.Log("error: issue with location id")
			http.Error(w, "issue with location id", http.StatusBadRequest)
			return
		}

		locationId := path[:8]
		path = strings.TrimPrefix(path, locationId)
		if !strings.HasPrefix(path, "/tickets") {
			h.Logger.Log("error: expecting tickets prefix")
			http.Error(w, "expecting tickets prefix", http.StatusBadRequest)
			return
		}

		h.mu.Lock()
		loc := h.Store.GetLocationById(locationId)
		h.mu.Unlock()
		if loc == nil {
			h.Logger.Log("error: missing location")
			http.NotFound(w, r)
			return
		}

		tq, err := models.NewTicketQuery(q)
		if err != nil {
			h.Logger.Log(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		tq.Location = loc
		if tq.IsOpen || tq.CloseEnd == 0 {
			h.Logger.Log("error: unexpected query param values")
			http.Error(w, "unexpected query param values", http.StatusBadRequest)
			return
		}

		h.mu.Lock()
		tq = h.Store.GetLocationTicketByTicketQuery(tq)
		h.mu.Unlock()

		resp := models.NewHalResp()
		resp.Embedded = tq.Result
		resp = PageResponse(resp, r, tq, h.Logger)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(resp)
		if err != nil {
			h.Logger.Log("error: issue encoding ticketList resp")
			http.Error(w, "issue encoding ticketList resp", http.StatusInternalServerError)
		}
		return
	case strings.HasPrefix(path, "/locations"):
		lq, err := models.NewLocationQuery(q)
		if err != nil {
			h.Logger.Log("error: unable to parse query")
			http.Error(w, "unable to parse query", http.StatusBadRequest)
			return
		}

		lq = h.Store.GetLocationsByLocationQuery(lq)
		resp := models.NewHalResp()
		resp.Embedded = lq.Result
		resp = PageResponse(resp, r, lq, h.Logger)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(resp)
		if err != nil {
			h.Logger.Log("error: issue encoding location resp")
			http.Error(w, "issue encoding location resp", http.StatusInternalServerError)
		}
		return
	default:
		http.NotFound(w, r)
		return
	}
}

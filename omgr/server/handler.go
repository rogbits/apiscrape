package server

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/oapiclient"
	"apiscrape/omgr/data"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Handler struct {
	Logger *loggr.Logger
	Store  *data.Store
}

func NewHandler(store *data.Store, logger *loggr.Logger) *Handler {
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

	if !strings.HasPrefix(u.Path, "/1.0") {
		http.Error(w, "missing version", http.StatusBadRequest)
		return
	}

	path := strings.TrimPrefix(u.Path, "/1.0")
	switch {
	case strings.HasPrefix(path, "/health"):
		_, err := w.Write([]byte("ok"))
		if err != nil {
			h.Logger.Log("error: writing health")
			http.Error(w, "error writing health", http.StatusInternalServerError)
		}
		return
	case strings.HasPrefix(path, "/_dequeue"):
		if len(h.Store.TicketFetchJobs) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.Store.Mu.Lock()
		fj := h.Store.TicketFetchJobs[0]
		h.Store.TicketFetchJobs = h.Store.TicketFetchJobs[1:]
		h.Store.JobsInFlight[fj.JobId] = fj
		h.Store.Mu.Unlock()

		h.Logger.Log("queue length", len(h.Store.TicketFetchJobs))

		encoder := json.NewEncoder(w)
		err := encoder.Encode(fj)
		if err != nil {
			h.Logger.Log("error writing job")
			http.Error(w, "error writing job", http.StatusInternalServerError)
		}
	case strings.HasPrefix(path, "/jobs"):
		fj := oapiclient.NewTicketFetchJob(0, "", "", 0, 0)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(fj)
		if err != nil {
			h.Logger.Log("error: bad job payload in request")
			http.Error(w, "bad job payload in request", http.StatusBadRequest)
		}

		h.Store.Mu.Lock()
		h.Store.ProcessJobResult(fj)
		delete(h.Store.JobsInFlight, fj.JobId)
		topItem, count := h.Store.GetTopItemsByLocationId(fj.LocationId)
		h.Logger.Log("latest top item for", fj.LocationName, ":", topItem, count)
		h.Store.Mu.Unlock()

		w.WriteHeader(http.StatusNoContent)
	case strings.HasPrefix(path, "/locations/"):
		path = strings.TrimPrefix(path, "/locations/")
		locId := path
		items, count := h.Store.GetTopItemsByLocationId(locId)
		resp := fmt.Sprintf("%s=%d\n", strings.Join(items, ":"), count)
		_, err := w.Write([]byte(resp))
		if err != nil {
			h.Logger.Log("error: couldn't write top items")
			http.Error(w, "couldn't write top items", http.StatusInternalServerError)
		}
	default:
		http.NotFound(w, r)
		return
	}
}

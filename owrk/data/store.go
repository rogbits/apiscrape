package data

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/models"
	"apiscrape/lib/oapiclient"
	"apiscrape/lib/omgrclient"
	"apiscrape/lib/tools"
	"time"
)

type OWrkStore struct {
	WorkerId                    string
	MillisecondsBetweenRequests int64
	Logger                      *loggr.Logger
	OMgrClient                  *omgrclient.OMgrClient
	OApiClient                  *oapiclient.OApiClient
}

func NewOWrkStore(logger *loggr.Logger,
	oMgrClient *omgrclient.OMgrClient, oApiClient *oapiclient.OApiClient) *OWrkStore {
	st := new(OWrkStore)
	st.WorkerId = tools.GenerateId(6)
	st.Logger = logger
	st.OMgrClient = oMgrClient
	st.OApiClient = oApiClient
	return st
}

func (store *OWrkStore) DequeueJobs() {
	backOff := time.Duration(store.MillisecondsBetweenRequests) * time.Millisecond
	increaseBackOff := func() {
		if backOff < 5*time.Minute {
			backOff += 10 * time.Second
		}
	}

	for {
		time.Sleep(backOff)

		// get job
		store.Logger.Log("------------")
		fj, err := store.OMgrClient.GetJob()
		if err != nil {
			store.Logger.Log("error fetching job", err)
			increaseBackOff()
			continue
		}
		if fj == nil {
			store.Logger.Log("no job available")
			increaseBackOff()
			continue
		}

		backOff = time.Duration(store.MillisecondsBetweenRequests) * time.Millisecond
		store.Logger.Log(
			"job received",
			"job id", fj.JobId,
			"loc id", fj.LocationId,
			"loc nm", fj.LocationName,
			"window start", time.Unix(fj.WindowStart, 0),
			"window end", time.Unix(fj.WindowEnd, 0),
		)

		err = store.ProcessJob(fj)
		if err != nil {
			store.Logger.Log("error processing job", err)
			increaseBackOff()
		} else {
			store.Logger.Log("job complete..")
		}
	}
}

func (store *OWrkStore) ProcessJob(fj *oapiclient.FetchTicketJob) error {
	tq := models.NewBlankTicketQuery()
	tq.IsOpen = false
	tq.CloseStart = fj.WindowStart
	tq.CloseEnd = fj.WindowEnd
	tq.PageStart = 1
	tq.PageLimit = 1
	tq.Location = models.NewLocation(fj.LocationName)
	tq.Location.Id = fj.LocationId
	tq.JobId = fj.JobId

	for {
		// fetch tickets
		store.Logger.Log(tq)
		tq, err := store.OApiClient.FetchTickets(tq)
		if err != nil {
			return err
		}
		for _, ticket := range tq.Result.Tickets {
			store.Logger.Log(
				"ticket rx",
				"num", ticket.TicketNumber,
				"job", tq.JobId,
				"next", tq.HalResp.Links.Next.Href,
			)
			for _, item := range ticket.Embedded.Items {
				fj.Items = append(fj.Items, item.Name)
				fj.ItemCounts = append(fj.ItemCounts, item.Quantity)
			}
		}
		if tq.HalResp.Links.Next.Href == "" {
			break
		}
		tq.PageStart, err = tq.HalResp.GetNextPageStart()
		if err != nil {
			return err
		}
		tq.PageLimit, err = tq.HalResp.GetNextPageLimit()
		if err != nil {
			return err
		}
	}
	store.Logger.Log("sending completed job to mgr with job id", fj.JobId)
	err := store.OMgrClient.UpdateJob(fj)
	if err != nil {
		return err
	}
	return nil
}

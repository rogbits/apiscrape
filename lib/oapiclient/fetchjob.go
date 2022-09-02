package oapiclient

type FetchTicketJob struct {
	JobId        int64  `json:"job_id"`
	LocationName string `json:"location_name"`
	LocationId   string `json:"location_id"`
	WindowStart  int64  `json:"window_start"`
	WindowEnd    int64  `json:"window_end"`

	Items      []string `json:"items"`
	ItemCounts []int64  `json:"item_counts"`
}

func NewTicketFetchJob(jobId int64, locId, locName string, windowStart, windowEnd int64) *FetchTicketJob {
	fj := new(FetchTicketJob)
	fj.JobId = jobId
	fj.LocationId = locId
	fj.LocationName = locName
	fj.WindowStart = windowStart
	fj.WindowEnd = windowEnd
	return fj
}

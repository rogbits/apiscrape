package server

import (
	"apiscrape/lib/loggr"
	"apiscrape/lib/models"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func PageResponse(resp *models.HalResp, r *http.Request, rq models.PaginatedQuery, logger *loggr.Logger) *models.HalResp {
	resp.Limit = rq.GetPageLimit()
	resp.Count = rq.GetPagedCount()

	url := r.URL.String()
	rtype := fmt.Sprintf("mock-api-%s", rq.GetResourceType())
	resp.Links.Self.Href = url
	resp.Links.Self.Type = rtype

	pageStart := rq.GetPageStart()
	pageLimit := rq.GetPageLimit()
	logger.Debug(
		"paging:",
		"start", pageStart,
		"limit", pageLimit,
		"count", rq.GetPagedCount(),
		"diff", pageStart-pageLimit,
		"sum", pageStart+pageLimit,
		"numRecords", rq.GetTotalNumRecords(),
	)

	// prev
	if pageStart > pageLimit {
		resp.Links.Prev.Href = SetPageParams(url, pageStart-pageLimit, pageLimit)
		resp.Links.Prev.Type = rtype
	}

	// next
	if pageStart-1+pageLimit < rq.GetTotalNumRecords() {
		resp.Links.Next.Href = SetPageParams(url, pageStart+pageLimit, pageLimit)
		resp.Links.Next.Type = rtype
	}

	return resp
}

func SetPageParams(url string, pageStart, pageLimit int64) string {
	r := regexp.MustCompile(`(&)?start=\d`)
	url = r.ReplaceAllString(url, "")

	r = regexp.MustCompile(`(&)?limit=\d`)
	url = r.ReplaceAllString(url, "")

	switch {
	case !strings.Contains(url, "?"):
		url += fmt.Sprintf(`?start=%d&limit=%d`, pageStart, pageLimit)
	case !strings.HasSuffix(url, "?"):
		url += fmt.Sprintf(`&start=%d&limit=%d`, pageStart, pageLimit)
	case strings.HasSuffix(url, "?"):
		url += fmt.Sprintf(`start=%d&limit=%d`, pageStart, pageLimit)
	}

	return url
}

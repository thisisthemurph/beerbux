package pagination

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	PageSize  int32
	PageToken string
}

// FromRequest creates a new Pagination struct from the given http.Request.
func FromRequest(r *http.Request) *Pagination {
	p := new(Pagination)

	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		ps, err := strconv.Atoi(pageSize)
		if err == nil {
			p.PageSize = int32(ps)
		}
	}

	if pageToken := r.URL.Query().Get("page_token"); pageToken != "" {
		p.PageToken = pageToken
	}

	return p
}

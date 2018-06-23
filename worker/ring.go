package worker

import (
	"container/ring"
	"strings"
	"time"
)

const (
	dateVar       = "{date}"
	isoDateFormat = "2006-01-02"
)

// QueryRing circulates over N queries
type QueryRing struct {
	Ring *ring.Ring
}

// NewQueryRing creates a new ring
func NewQueryRing(items []string) *QueryRing {
	r := &QueryRing{
		Ring: ring.New(len(items)),
	}

	for _, item := range items {
		r.Ring.Value = item
		r.Ring = r.Ring.Next()
	}

	return r
}

// Next entry from the query ring
func (r *QueryRing) Next() string {
	v := r.Ring.Value
	r.Ring = r.Ring.Next()
	date := time.Now().AddDate(0, 0, -1).UTC().Format(isoDateFormat)
	return strings.Replace(v.(string), dateVar, date, -1)
}

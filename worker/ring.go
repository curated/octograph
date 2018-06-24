package worker

import (
	"container/ring"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	dateVar       = "{date}"
	isoDateFormat = "2006-01-02"
	reactionsVar  = "{reactions}"
	reactionsMin  = 25
	reactionsMax  = 150
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
	return r.replaceVars(v.(string))
}

// Rollback to previous entry
func (r *QueryRing) Rollback() {
	r.Ring = r.Ring.Prev()
}

func (r *QueryRing) replaceVars(s string) string {
	date := time.Now().AddDate(0, 0, -1).UTC().Format(isoDateFormat)
	s = strings.Replace(s, dateVar, date, -1)

	rand.Seed(time.Now().UnixNano())
	reactions := reactionsMin + rand.Intn(reactionsMax-reactionsMin)
	return strings.Replace(s, reactionsVar, strconv.Itoa(reactions), -1)
}

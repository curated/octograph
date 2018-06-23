package worker_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/curated/octograph/worker"
	"github.com/stretchr/testify/assert"
)

func TestQueryRing(t *testing.T) {
	r := worker.NewQueryRing([]string{
		"foo: {date}",
		"bar",
		"zip",
	})

	date := time.Now().AddDate(0, 0, -1).UTC().Format("2006-01-02")
	foo := fmt.Sprintf("foo: %s", date)

	assert.Equal(t, foo, r.Next())
	assert.Equal(t, "bar", r.Next())
	assert.Equal(t, "zip", r.Next())
	assert.Equal(t, foo, r.Next())
}

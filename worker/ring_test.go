package worker_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/curated/octograph/worker"
	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	r := worker.NewQueryRing([]string{
		"foo",
		"bar",
		"zip",
	})

	assert.Equal(t, "foo", r.Next())
	assert.Equal(t, "bar", r.Next())
	assert.Equal(t, "zip", r.Next())
	assert.Equal(t, "foo", r.Next())
}

func TestRollback(t *testing.T) {
	r := worker.NewQueryRing([]string{
		"foo",
		"bar",
	})

	assert.Equal(t, "foo", r.Next())
	assert.Equal(t, "bar", r.Next())

	r.Rollback()
	assert.Equal(t, "bar", r.Next())
}

func TestDate(t *testing.T) {
	r := worker.NewQueryRing([]string{
		"foo: {date}",
		"bar",
	})

	date := time.Now().AddDate(0, 0, -1).UTC().Format("2006-01-02")
	foo := fmt.Sprintf("foo: %s", date)

	assert.Equal(t, foo, r.Next())
	assert.Equal(t, "bar", r.Next())
	assert.Equal(t, foo, r.Next())
}

func TestReactions(t *testing.T) {
	r := worker.NewQueryRing([]string{
		"foo: {reactions}",
		"bar",
	})

	foo := r.Next()
	assert.Equal(t, "foo: ", foo[:5])
	reactions, err := strconv.Atoi(foo[5:])
	assert.Nil(t, err)
	assert.True(
		t,
		reactions >= 25 && reactions <= 150,
		"Reactions should be in range 25:150",
	)

	assert.Equal(t, "bar", r.Next())

	foo = r.Next()
	assert.Equal(t, "foo: ", foo[:5])
	reactions, err = strconv.Atoi(foo[5:])
	assert.Nil(t, err)
	assert.True(
		t,
		reactions >= 25 && reactions <= 150,
		"Reactions should be in range 25:150",
	)
}

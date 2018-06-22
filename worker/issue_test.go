package worker_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/indexer"
	"github.com/curated/octograph/worker"
	"github.com/stretchr/testify/assert"
)

var c = config.New("../")
var issueWorker = worker.NewIssueWorker(c)
var idx = indexer.New(c)
var index = "test"

func TestMain(m *testing.M) {
	exists, err := idx.Client.IndexExists(index).Do(idx.Context)
	if err != nil {
		panic(err)
	}

	if exists {
		res, err := idx.Client.DeleteIndex(index).Do(idx.Context)
		if err != nil || !res.Acknowledged {
			panic(fmt.Sprintf("Delete acknowledged: %v, error: %v", res.Acknowledged, err))
		}
	}

	res, err := idx.Client.CreateIndex(index).Do(idx.Context)
	if err != nil {
		panic(fmt.Sprintf("Create acknowledged: %v, error: %v", res.Acknowledged, err))
	}

	os.Exit(m.Run())
}

func TestProcess(t *testing.T) {
	query := "reactions:>3000"

	sr, err := idx.Client.Search().
		Index(index).
		From(0).
		Size(10).
		Do(idx.Context)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), sr.TotalHits())

	issueWorker.Process(query)
	assert.Nil(t, err)

	_, err = idx.Client.Flush(index).Do(idx.Context)
	assert.Nil(t, err)

	sr, err = idx.Client.Search().
		Index(index).
		From(0).
		Size(10).
		Do(idx.Context)

	assert.Nil(t, err)
	assert.True(t, sr.TotalHits() > 0)
}

package worker_test

import (
	"os"
	"testing"

	"github.com/curated/octograph/indexer"
	"github.com/curated/octograph/worker"
	"github.com/stretchr/testify/assert"
)

var issueWorker = worker.NewIssueWorker()
var idx = indexer.New()
var index = "issue"

func TestMain(m *testing.M) {
	exists, err := idx.Client.IndexExists(index).Do(idx.Context)
	if err != nil {
		panic(err)
	}
	if exists {
		_, err := idx.Client.DeleteIndex(index).Do(idx.Context)
		if err != nil {
			panic(err)
		}
	}
	_, err = idx.Client.CreateIndex(index).Do(idx.Context)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestProcess(t *testing.T) {
	sr, err := idx.Client.Search().
		Index(index).
		From(0).
		Size(50).
		Do(idx.Context)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), sr.TotalHits())

	issueWorker.Process()
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

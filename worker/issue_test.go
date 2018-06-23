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

func TestMain(m *testing.M) {
	exists, err := idx.Client.IndexExists(c.Issue.Index).Do(idx.Context)
	if err != nil {
		panic(err)
	}

	if exists {
		res, err := idx.Client.DeleteIndex(c.Issue.Index).Do(idx.Context)
		if err != nil {
			panic(err)
		}
		if !res.Acknowledged {
			panic(fmt.Sprintf("Index deletion was not acknowledged: %+v", res))
		}
	}

	res, err := idx.Client.CreateIndex(c.Issue.Index).Do(idx.Context)
	if err != nil {
		panic(err)
	}
	if !res.Acknowledged {
		panic(fmt.Sprintf("Index creation was not acknowledged: %+v", res))
	}

	os.Exit(m.Run())
}

func TestIndex(t *testing.T) {
	sr, err := idx.Client.Search().
		Index(c.Issue.Index).
		From(0).
		Size(10).
		Do(idx.Context)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), sr.TotalHits())

	issueWorker.Index()
	assert.Nil(t, err)

	_, err = idx.Client.Flush(c.Issue.Index).Do(idx.Context)
	assert.Nil(t, err)

	sr, err = idx.Client.Search().
		Index(c.Issue.Index).
		From(0).
		Size(10).
		Do(idx.Context)

	assert.Nil(t, err)
	assert.True(t, sr.TotalHits() > 0)
}

func TestDelete(t *testing.T) {
	err := issueWorker.Delete()
	assert.Nil(t, err)

	exists, err := idx.Client.IndexExists(c.Issue.Index).Do(idx.Context)

	assert.Nil(t, err)
	assert.False(t, exists)
}

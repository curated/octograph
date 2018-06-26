package indexer_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/indexer"
	"github.com/stretchr/testify/assert"
)

var c = config.New("../")
var idx = indexer.New(c)
var issueType = "issue"

func TestMain(m *testing.M) {
	c.Issue.Index = "indexer_test"

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

	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	err := idx.Create(c.Issue.Index, "{}")
	assert.Nil(t, err)

	exists, err := idx.Client.IndexExists(c.Issue.Index).Do(idx.Context)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestEnsure(t *testing.T) {
	res, err := idx.Client.DeleteIndex(c.Issue.Index).Do(idx.Context)
	assert.Nil(t, err)
	assert.True(t, res.Acknowledged)

	for i := 1; i <= 2; i++ {
		err = idx.Ensure(c.Issue.Index, "{}")
		assert.Nil(t, err)

		exists, err := idx.Client.IndexExists(c.Issue.Index).Do(idx.Context)
		assert.Nil(t, err)
		assert.True(t, exists)
	}
}

func TestDelete(t *testing.T) {
	err := idx.Delete(c.Issue.Index)
	assert.Nil(t, err)

	exists, err := idx.Client.IndexExists(c.Issue.Index).Do(idx.Context)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestIndex(t *testing.T) {
	res, err := idx.Client.CreateIndex(c.Issue.Index).Do(idx.Context)
	assert.Nil(t, err)
	assert.True(t, res.Acknowledged)

	sr, err := idx.Client.Search().
		Index(c.Issue.Index).
		From(0).
		Size(10).
		Do(idx.Context)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), sr.TotalHits())

	err = idx.Index(c.Issue.Index, issueType, "id", "{}")
	assert.Nil(t, err)

	_, err = idx.Client.Flush(c.Issue.Index).Do(idx.Context)
	assert.Nil(t, err)

	sr, err = idx.Client.Search().
		Index(c.Issue.Index).
		From(0).
		Size(10).
		Do(idx.Context)

	assert.Nil(t, err)
	assert.Equal(t, int64(1), sr.TotalHits())
}

func TestGet(t *testing.T) {
	err := idx.Index(c.Issue.Index, issueType, "id", "{\"foo\": \"bar\"}")
	assert.Nil(t, err)

	b, err := idx.Get(c.Issue.Index, issueType, "id")
	var source map[string]interface{}
	err = json.Unmarshal(b, &source)

	assert.Nil(t, err)
	assert.Equal(t, "bar", source["foo"])
}

func TestGetMapping(t *testing.T) {
	s, err := idx.GetMapping("issue.json")
	assert.Nil(t, err)

	var m map[string]interface{}
	err = json.Unmarshal([]byte(s), &m)
	assert.Nil(t, err)
	assert.True(t, len(s) > 0)
}

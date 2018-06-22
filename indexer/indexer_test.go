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

var idx = indexer.New(config.New("../"))
var index = "test"

func TestMain(m *testing.M) {
	exists, err := idx.Client.IndexExists(index).Do(idx.Context)
	if err != nil {
		panic(err)
	}

	if exists {
		res, err := idx.Client.DeleteIndex(index).Do(idx.Context)
		if err != nil {
			panic(fmt.Sprintf("Delete acknowledged: %v, error: %v", res.Acknowledged, err))
		}
	}

	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	err := idx.Create(index, "{}")
	assert.Nil(t, err)

	exists, err := idx.Client.IndexExists(index).Do(idx.Context)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestEnsure(t *testing.T) {
	res, err := idx.Client.DeleteIndex(index).Do(idx.Context)
	assert.True(t, res.Acknowledged)
	assert.Nil(t, err)

	for i := 1; i <= 2; i++ {
		err = idx.Ensure(index, "{}")
		assert.Nil(t, err)

		exists, err := idx.Client.IndexExists(index).Do(idx.Context)
		assert.Nil(t, err)
		assert.True(t, exists)
	}
}

func TestDelete(t *testing.T) {
	err := idx.Delete(index)
	assert.Nil(t, err)

	exists, err := idx.Client.IndexExists(index).Do(idx.Context)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestIndex(t *testing.T) {
	res, err := idx.Client.CreateIndex(index).Do(idx.Context)
	assert.True(t, res.Acknowledged)
	assert.Nil(t, err)

	sr, err := idx.Client.Search().
		Index(index).
		From(0).
		Size(10).
		Do(idx.Context)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), sr.TotalHits())

	err = idx.Index(index, index, "id", "{}")
	assert.Nil(t, err)

	_, err = idx.Client.Flush(index).Do(idx.Context)
	assert.Nil(t, err)

	sr, err = idx.Client.Search().
		Index(index).
		From(0).
		Size(10).
		Do(idx.Context)

	assert.Nil(t, err)
	assert.Equal(t, int64(1), sr.TotalHits())
}

func TestGetMapping(t *testing.T) {
	s, err := idx.GetMapping("issue.json")
	assert.Nil(t, err)

	var m map[string]interface{}
	err = json.Unmarshal([]byte(s), &m)
	assert.Nil(t, err)

	assert.True(t, len(s) > 0)
}

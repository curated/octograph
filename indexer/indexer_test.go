package indexer_test

import (
	"os"
	"testing"

	"github.com/curated/octograph/indexer"
	"github.com/stretchr/testify/assert"
)

var idx = indexer.New()
var index = "test"

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
	_, err := idx.Client.DeleteIndex(index).Do(idx.Context)
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
	sr, err := idx.Client.Search().
		Index(index).
		From(0).
		Size(10).
		Do(idx.Context)
	assert.NotNil(t, err)

	_, err = idx.Client.CreateIndex(index).Do(idx.Context)
	assert.Nil(t, err)

	sr, err = idx.Client.Search().
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

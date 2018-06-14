package graph_test

import (
	"encoding/json"
	"testing"

	"github.com/curated/octograph/graph"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	g := graph.New()

	query := []byte(`
		query {
			viewer {
				isViewer
				login
			}
		}
	`)

	b, err := g.Fetch(query, nil)
	assert.Nil(t, err)

	type Response struct {
		Data struct {
			Viewer struct {
				Login string
			}
		}
	}

	var res Response
	err = json.Unmarshal(b, &res)
	assert.Nil(t, err)

	assert.True(t, len(res.Data.Viewer.Login) > 0)
}

func TestFetchWithVariable(t *testing.T) {
	g := graph.New()

	query := []byte(`
		query Organization($login: String!) {
			organization(login: $login) {
				url
			}
		}
	`)

	variables := map[string]interface{}{
		"login": "curated",
	}

	b, err := g.Fetch(query, variables)
	assert.Nil(t, err)

	type Response struct {
		Data struct {
			Organization struct {
				URL string
			}
		}
	}

	var res Response
	err = json.Unmarshal(b, &res)
	assert.Nil(t, err)

	assert.Equal(t, "https://github.com/curated", res.Data.Organization.URL)
}

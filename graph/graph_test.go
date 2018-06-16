package graph_test

import (
	"encoding/json"
	"testing"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/graph"

	"github.com/stretchr/testify/assert"
)

var g = graph.New(config.New("../"))

func TestFetch(t *testing.T) {
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

func TestFetchIssues(t *testing.T) {
	query := "reactions:>1000"
	issues, err := g.FetchIssues(query, nil)
	assert.Nil(t, err)

	assert.True(t, issues.Data.Search.IssueCount > 0)
	assert.True(t, len(issues.Data.Search.Edges) > 0)

	issuesAfter, err := g.FetchIssues(query, &issues.Data.Search.PageInfo.EndCursor)
	assert.Nil(t, err)

	assert.True(t, issuesAfter.Data.Search.IssueCount > 0)
	assert.True(t, len(issuesAfter.Data.Search.Edges) > 0)
}

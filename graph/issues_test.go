package graph_test

import (
	"testing"

	"github.com/curated/octograph/graph"

	"github.com/stretchr/testify/assert"
)

func TestFetchIssues(t *testing.T) {
	g := graph.New()

	issues, err := g.FetchIssues(nil)
	assert.Nil(t, err)

	assert.True(t, issues.Data.Search.IssueCount > 0)
	assert.True(t, len(issues.Data.Search.Edges) > 0)

	issuesAfter, err := g.FetchIssues(&issues.Data.Search.PageInfo.EndCursor)
	assert.Nil(t, err)

	assert.True(t, issuesAfter.Data.Search.IssueCount > 0)
	assert.True(t, len(issuesAfter.Data.Search.Edges) > 0)
}

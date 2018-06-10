package graph

import (
	"io/ioutil"
)

const topIssuesGQL = "./graph/TopIssues.gql"

// FetchTopIssues after optional end cursor
func (g *Graph) FetchTopIssues(endCursor *string) ([]byte, error) {
	g.Logger.Printf("Fetching top issues after cursor: %v", endCursor)

	b, err := ioutil.ReadFile(topIssuesGQL)
	if err != nil {
		g.Logger.Printf("Failed fetching top issues: %v", err)
		return []byte{}, err
	}

	return g.fetch(b, map[string]interface{}{
		"after": endCursor,
	})
}

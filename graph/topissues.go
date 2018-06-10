package graph

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

const topIssuesGQL = "./graph/TopIssues.gql"

// FetchTopIssues after optional end cursor
func (g *Graph) FetchTopIssues(endCursor *string) (*TopIssues, error) {
	g.Logger.Printf("Fetching top issues after cursor: %v", endCursor)

	query, err := ioutil.ReadFile(topIssuesGQL)
	if err != nil {
		g.Logger.Printf("Failed reading '%s' with error: %v", topIssuesGQL, err)
		return nil, err
	}

	res, err := g.fetch(query, map[string]interface{}{
		"after": endCursor,
	})
	if err != nil {
		g.Logger.Printf("Failed fetching top issues: %v", err)
		return nil, err
	}

	var ti TopIssues
	err = json.Unmarshal(res, &ti)
	if err != nil {
		g.Logger.Printf("Failed parsing top issues: %v", err)
		return nil, err
	}

	return &ti, nil
}

// TopIssues response structure from GraphQL endpoint
type TopIssues struct {
	Data struct {
		Search struct {
			IssueCount int

			PageInfo struct {
				EndCursor string
			}

			Edges []struct {
				Node struct {
					ID        string
					URL       string
					Number    int
					Title     string
					BodyText  string
					State     string
					CreatedAt time.Time
					UpdatedAt time.Time

					ReactionGroups []struct {
						Content string

						Users struct {
							TotalCount int
						}
					}

					Repository struct {
						ID   string
						URL  string
						Name string

						PrimaryLanguage struct {
							Name string
						}

						Forks struct {
							TotalCount int
						}

						Stargazers struct {
							TotalCount int
						}

						Owner struct {
							ID        string
							URL       string
							Name      string
							AvatarURL string
						}
					}

					Author struct {
						ID        string
						URL       string
						Login     string
						AvatarURL string
					}
				}
			}
		}
	}
}

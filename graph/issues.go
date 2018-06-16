package graph

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/curated/octograph/config"
	"github.com/golang/glog"
)

// Issues response structure from GraphQL endpoint
type Issues struct {
	Data struct {
		Search struct {
			IssueCount int

			PageInfo struct {
				EndCursor string
			}

			Edges []struct {
				Node Issue
			}
		}
	}
}

// Issue node structure from GraphQL endpoint
type Issue struct {
	ID        string
	URL       string
	Number    int
	Title     string
	BodyText  string
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time

	ReactionGroups []ReactionGroup

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
			Login     string
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

// ReactionGroup node structure from GraphQL
type ReactionGroup struct {
	Content string

	Users struct {
		TotalCount int
	}
}

// FetchIssues by query, after optional end cursor
func (g *Graph) FetchIssues(query string, endCursor *string) (*Issues, error) {
	issuesGQL := config.GetPath("graph/issues.gql")
	b, err := ioutil.ReadFile(issuesGQL)

	if err != nil {
		glog.Errorf("Failed reading '%s' with error: %v", issuesGQL, err)
		return nil, err
	}

	res, err := g.Fetch(b, map[string]interface{}{
		"query": query,
		"after": endCursor,
	})

	if err != nil {
		glog.Errorf("Failed fetching issues: %v", err)
		return nil, err
	}

	var issues Issues
	err = json.Unmarshal(res, &issues)

	if err != nil {
		glog.Errorf("Failed parsing issues: %v\n%s", err, string(res))
		return nil, err
	}

	return &issues, nil
}

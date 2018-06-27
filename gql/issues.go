package gql

import (
	"time"
)

// Issues response structure from GraphQL
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

// Issue node structure from GraphQL
type Issue struct {
	ID        string
	Number    int
	Title     string
	Body      string
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time

	ReactionGroups []ReactionGroup

	Author struct {
		Login string
	}

	Repository struct {
		URL string

		PrimaryLanguage struct {
			Name string
		}

		Forks struct {
			TotalCount int
		}

		Stargazers struct {
			TotalCount int
		}
	}
}

// ReactionGroup node structure from GraphQL
type ReactionGroup struct {
	Content string

	Users struct {
		TotalCount int
	}
}

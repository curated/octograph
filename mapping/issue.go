package mapping

import (
	"time"
)

// Issue serialization structure for indexing
type Issue struct {
	Number         int       `json:"number"`
	Title          string    `json:"title"`
	BodyText       string    `json:"bodyText"`
	State          string    `json:"state"`
	ThumbsUp       int       `json:"thumbsUp"`
	ThumbsDown     int       `json:"thumbsDown"`
	Laugh          int       `json:"laugh"`
	Hooray         int       `json:"hooray"`
	Confused       int       `json:"confused"`
	Heart          int       `json:"heart"`
	AuthorLogin    string    `json:"authorLogin"`
	RepoOwnerLogin string    `json:"repoOwnerLogin"`
	RepoName       string    `json:"repoName"`
	RepoLanguage   string    `json:"repoLanguage"`
	RepoForks      int       `json:"repoForks"`
	RepoStargazers int       `json:"repoStargazers"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

package mapping

import (
	"time"
)

// Issue serialization structure for indexing
type Issue struct {
	ID                    string    `json:"id"`
	URL                   string    `json:"url"`
	Number                int       `json:"number"`
	Title                 string    `json:"title"`
	BodyText              string    `json:"bodyText"`
	State                 string    `json:"state"`
	ThumbsUp              int       `json:"thumbsUp"`
	ThumbsDown            int       `json:"thumbsDown"`
	Laugh                 int       `json:"laugh"`
	Hooray                int       `json:"hooray"`
	Confused              int       `json:"confused"`
	Heart                 int       `json:"heart"`
	AuthorLogin           string    `json:"authorLogin"`
	AuthorLoginSuggest    string    `json:"authorLoginSuggest"`
	RepoOwnerLogin        string    `json:"repoOwnerLogin"`
	RepoOwnerLoginSuggest string    `json:"repoOwnerLoginSuggest"`
	RepoName              string    `json:"repoName"`
	RepoNameSuggest       string    `json:"repoNameSuggest"`
	RepoLanguage          string    `json:"repoLanguage"`
	RepoLanguageSuggest   string    `json:"repoLanguageSuggest"`
	RepoForks             int       `json:"repoForks"`
	RepoStargazers        int       `json:"repoStargazers"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

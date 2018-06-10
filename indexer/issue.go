package indexer

import (
	"io/ioutil"
	"time"
)

const (
	issueMappingJSON = "./indexer/issue.json"
	issueIndexName   = "issue"
	issueIndexType   = "issue"
)

// Issue serialization structure for indexing
type Issue struct {
	ID              string    `json:"id"`
	URL             string    `json:"url"`
	Number          int       `json:"number"`
	Title           string    `json:"title"`
	BodyText        string    `json:"bodyText"`
	State           string    `json:"state"`
	ThumbsUp        int       `json:"thumbsUp"`
	ThumbsDown      int       `json:"thumbsDown"`
	Laugh           int       `json:"laugh"`
	Hooray          int       `json:"hooray"`
	Confused        int       `json:"confused"`
	Heart           int       `json:"heart"`
	AuthorID        string    `json:"authorId"`
	AuthorURL       string    `json:"authorUrl"`
	AuthorLogin     string    `json:"authorLogin"`
	AuthorAvatar    string    `json:"authorAvatar"`
	RepoID          string    `json:"repoId"`
	RepoURL         string    `json:"repoUrl"`
	RepoName        string    `json:"repoName"`
	RepoLanguage    string    `json:"repoLanguage"`
	RepoForks       int       `json:"repoForks"`
	RepoStargazers  int       `json:"repoStargazers"`
	RepoOwnerID     string    `json:"repoOwnerId"`
	RepoOwnerURL    string    `json:"repoOwnerUrl"`
	RepoOwnerLogin  string    `json:"repoOwnerLogin"`
	RepoOwnerAvatar string    `json:"repoOwnerAvatar"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func issueMapping() (string, error) {
	b, err := ioutil.ReadFile(issueMappingJSON)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

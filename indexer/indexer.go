package indexer

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/curated/octograph/config"
	"github.com/golang/glog"
	"github.com/olivere/elastic"
)

const (
	// IssueIndex in Elastic
	IssueIndex = "issue"

	// IssueType in Elastic
	IssueType = "issue"

	elasticScheme   = "https"
	elasticSniffing = false
)

// Indexer client
type Indexer struct {
	Context context.Context
	Client  *elastic.Client
	Config  *config.Config
}

// Issue serialization structure for indexing
type Issue struct {
	ID             string    `json:"id"`
	URL            string    `json:"url"`
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

// New create a new indexer
func New(c *config.Config) *Indexer {
	cli, err := elastic.NewClient(
		elastic.SetURL(c.Elastic.URL),
		elastic.SetBasicAuth(c.Elastic.Username, c.Elastic.Password),
		elastic.SetScheme(elasticScheme),
		elastic.SetSniff(elasticSniffing),
	)

	if err != nil {
		glog.Fatalf("Failed creating indexer: %v", err)
	}

	return &Indexer{
		Context: context.Background(),
		Client:  cli,
		Config:  c,
	}
}

// Create index
func (i *Indexer) Create(index, mapping string) error {
	res, err := i.Client.CreateIndex(index).BodyString(mapping).Do(i.Context)

	if err != nil {
		glog.Errorf("Failed creating index: %v", err)
		return err
	}

	if !res.Acknowledged {
		return fmt.Errorf("Index creation was not acknowledged: %v", res)
	}

	return nil
}

// Ensure index exists
func (i *Indexer) Ensure(index, mapping string) error {
	exists, err := i.Client.IndexExists(index).Do(i.Context)

	if err != nil {
		glog.Errorf("Failed checking if index exists: %v", err)
		return err
	}

	if !exists {
		return i.Create(index, mapping)
	}

	return nil
}

// Delete index
func (i *Indexer) Delete(index string) error {
	res, err := i.Client.DeleteIndex(index).Do(i.Context)

	if err != nil {
		glog.Errorf("Failed deleting index: %v", err)
		return err
	}

	if !res.Acknowledged {
		return fmt.Errorf("Index deletion was not acknowledged: %v", res)
	}

	return nil
}

// Index a document
func (i *Indexer) Index(name, typ, id string, body interface{}) error {
	_, err := i.Client.Index().
		Index(name).
		Type(typ).
		Id(id).
		BodyJson(body).
		Do(i.Context)

	if err != nil {
		glog.Errorf("Failed indexing document: %v", err)
		return err
	}

	return nil
}

// IssueMapping in Elastic
func (i *Indexer) IssueMapping() (string, error) {
	mappingJSON := i.Config.GetPath("indexer/issue_mapping.json")
	b, err := ioutil.ReadFile(mappingJSON)

	if err != nil {
		glog.Errorf("Failed reading '%s' with error: %v", mappingJSON, err)
		return "", err
	}

	return string(b), nil
}

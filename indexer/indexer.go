package indexer

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/curated/octograph/config"
	"github.com/golang/glog"
	"github.com/olivere/elastic"
)

const (
	// ElasticErrorNotFound from get query
	ElasticErrorNotFound = "Error 404 (Not Found)"

	elasticScheme   = "https"
	elasticSniffing = false
)

// Indexer client
type Indexer struct {
	Context context.Context
	Client  *elastic.Client
	Config  *config.Config
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

// Get document source
func (i *Indexer) Get(index, typ, id string) ([]byte, error) {
	res, err := i.Client.Get().
		Index(index).
		Type(typ).
		Id(id).
		Do(i.Context)

	if err != nil {
		if err.Error() == ElasticErrorNotFound {
			return nil, err
		}
		glog.Errorf("Failed getting document: %v", err)
		return nil, err
	}

	return *res.Source, err
}

// GetMapping for index
func (i *Indexer) GetMapping(filename string) (string, error) {
	f := i.Config.GetPath(fmt.Sprintf("mapping/%s", filename))
	b, err := ioutil.ReadFile(f)

	if err != nil {
		glog.Errorf("Failed reading '%s' with error: %v", f, err)
		return "", err
	}

	return string(b), nil
}

package indexer

import (
	"context"
	"fmt"

	"github.com/curated/octograph/config"
	"github.com/golang/glog"
	"github.com/olivere/elastic"
)

// New create a new indexer
func New() *Indexer {
	cfg := config.New()

	cli, err := elastic.NewClient(
		elastic.SetURL(cfg.Elastic.URL),
		elastic.SetBasicAuth(cfg.Elastic.Username, cfg.Elastic.Password),
		elastic.SetScheme("https"),
		elastic.SetSniff(false),
	)

	if err != nil {
		glog.Fatalf("Failed creating indexer: %v", err)
	}

	return &Indexer{
		Context: context.Background(),
		Client:  cli,
		Config:  cfg,
	}
}

// Indexer client
type Indexer struct {
	Context context.Context
	Client  *elastic.Client
	Config  *config.Config
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
		return fmt.Errorf("Index removal was not acknowledged: %v", res)
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

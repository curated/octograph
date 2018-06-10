package indexer

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/logger"
	"github.com/olivere/elastic"
)

// New create a new indexer
func New() *Indexer {
	cfg := config.Load()
	lg := logger.New()

	cli, err := elastic.NewClient(
		elastic.SetURL(cfg.Elastic.URL),
		elastic.SetBasicAuth(cfg.Elastic.Username, cfg.Elastic.Password),
		elastic.SetScheme("https"),
		elastic.SetSniff(false),
	)

	if err != nil {
		lg.Printf("Failed creating indexer: %v", err)
		os.Exit(1)
	}

	return &Indexer{
		Context: context.Background(),
		Client:  cli,
		Config:  cfg,
		Logger:  lg,
	}
}

// Indexer client
type Indexer struct {
	Context context.Context
	Client  *elastic.Client
	Config  *config.Config
	Logger  *log.Logger
}

// Create index
func (i *Indexer) Create(index, mapping string) error {
	res, err := i.Client.CreateIndex(index).BodyString(mapping).Do(i.Context)
	if err != nil {
		i.Logger.Printf("Failed creating index: %v", err)
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
		i.Logger.Printf("Failed checking if index exists: %v", err)
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
		i.Logger.Printf("Failed deleting index: %v", err)
		return err
	}
	if !res.Acknowledged {
		return fmt.Errorf("Index removal was not acknowledged: %v", res)
	}
	return nil
}

// Recreate index
func (i *Indexer) Recreate(index, mapping string) error {
	err := i.Delete(index)
	if err != nil {
		return err
	}
	return i.Create(index, mapping)
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
		i.Logger.Printf("Failed indexing document: %v", err)
		return err
	}
	return nil
}

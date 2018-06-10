package worker

import (
	"log"

	"github.com/curated/octograph/graph"
	"github.com/curated/octograph/logger"
)

// New creates a new worker
func New() *Worker {
	return &Worker{
		Graph:  graph.New(),
		Logger: logger.New(),
	}
}

// Worker struct
type Worker struct {
	Graph  *graph.Graph
	Logger *log.Logger
}

// Process indexing of top issues
func (w *Worker) Process() error {
	_, err := w.Graph.FetchTopIssues(nil)
	if err != nil {
		return err
	}

	return nil
}

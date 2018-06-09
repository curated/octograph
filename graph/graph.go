package graph

import (
	"io/ioutil"
	"log"
	"net/http"
)

// New creates a new graph client
func New(logger *log.Logger) *Graph {
	return &Graph{
		Logger: logger,
	}
}

// Graph client
type Graph struct {
	Logger *log.Logger
}

// Fetch issues
func (g *Graph) Fetch() (int, error) {
	g.Logger.Printf("Fetching issues")
	response, err := http.Get("https://api.github.com/graphql")
	if err != nil {
		g.Logger.Printf("Failed fetching issues: %v", err)
		return 0, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		g.Logger.Printf("Failed reading body: %v", err)
		return 0, err
	}

	return len(body), nil
}

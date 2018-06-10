package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/logger"
)

const endpoint = "https://api.github.com/graphql"

// New creates a new graph client
func New() *Graph {
	return &Graph{
		Config: config.Load(),
		Logger: logger.New(),
	}
}

// Graph client
type Graph struct {
	Config *config.Config
	Logger *log.Logger
}

// ReqBody structure for GraphQL endpoint
type ReqBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (g *Graph) fetch(query []byte, variables map[string]interface{}) ([]byte, error) {
	g.Logger.Printf("Query:\n%s", string(query))
	g.Logger.Printf("Variables: %v", variables)

	reqBody, err := json.Marshal(&ReqBody{
		Query:     string(query),
		Variables: variables,
	})
	if err != nil {
		g.Logger.Printf("Failed parsing request body: %v", err)
		return []byte{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		g.Logger.Printf("Failed creating request: %v", err)
		return []byte{}, err
	}

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", g.Config.GitHub.Token),
	)

	cli := http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		g.Logger.Printf("Failed processing request: %v", err)
		return []byte{}, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		g.Logger.Printf("Failed reading response body: %v", err)
		return []byte{}, err
	}

	g.Logger.Printf("Response headers: %v", res.Header)
	g.Logger.Printf("Response body:\n%s", string(resBody))
	return resBody, nil
}

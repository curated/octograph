package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/curated/octograph/config"
	"github.com/golang/glog"
)

const endpoint = "https://api.github.com/graphql"

// New creates a new graph client
func New() *Graph {
	return &Graph{
		Client: &http.Client{},
		Config: config.New(),
	}
}

// Graph client
type Graph struct {
	Client *http.Client
	Config *config.Config
}

// ReqBody structure for GraphQL endpoint
type ReqBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// Fetch GraphQL data using query and variables
func (g *Graph) Fetch(query []byte, variables map[string]interface{}) ([]byte, error) {
	reqBody, err := json.Marshal(&ReqBody{
		Query:     string(query),
		Variables: variables,
	})
	if err != nil {
		glog.Errorf("Failed parsing request body: %v\n%s\n%v", err, string(query), variables)
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		glog.Errorf("Failed creating request: %v", err)
		return nil, err
	}

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", g.Config.GitHub.Token),
	)

	res, err := g.Client.Do(req)
	if err != nil {
		glog.Errorf("Failed processing request: %v", err)
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		glog.Errorf("Failed reading response body: %v", err)
		return nil, err
	}

	return resBody, nil
}

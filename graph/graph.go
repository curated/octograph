package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/gql"
	"github.com/golang/glog"
)

const endpoint = "https://api.github.com/graphql"

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

// New creates a new graph client
func New(c *config.Config) *Graph {
	return &Graph{
		Client: &http.Client{},
		Config: c,
	}
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

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"Request failed: %d: %s: %s",
			res.StatusCode,
			res.Status,
			string(resBody),
		)
	}

	return resBody, nil
}

// FetchIssues by query, after optional end cursor
func (g *Graph) FetchIssues(query string, endCursor *string) (*gql.Issues, error) {
	issuesGQL := g.Config.GetPath("gql/issues.gql")
	b, err := ioutil.ReadFile(issuesGQL)

	if err != nil {
		glog.Errorf("Failed reading '%s' with error: %v", issuesGQL, err)
		return nil, err
	}

	res, err := g.Fetch(b, map[string]interface{}{
		"query": query,
		"after": endCursor,
	})

	if err != nil {
		glog.Errorf("Failed fetching issues: %v", err)
		return nil, err
	}

	var issues gql.Issues
	err = json.Unmarshal(res, &issues)

	if err != nil {
		glog.Errorf("Failed parsing issues: %v\n%s", err, string(res))
		return nil, err
	}

	return &issues, nil
}

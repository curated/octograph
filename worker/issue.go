package worker

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/gql"
	"github.com/curated/octograph/graph"
	"github.com/curated/octograph/indexer"
	"github.com/curated/octograph/mapping"
	"github.com/golang/glog"
	"github.com/google/go-cmp/cmp"
)

const (
	issueType = "issue"

	reactionThumbsUp   = "THUMBS_UP"
	reactionThumbsDown = "THUMBS_DOWN"
	reactionLaugh      = "LAUGH"
	reactionHooray     = "HOORAY"
	reactionConfused   = "CONFUSED"
	reactionHeart      = "HEART"

	elasticErrorNotFound = "Error 404 (Not Found)"
	missingValue         = "?"
)

// IssueWorker struct
type IssueWorker struct {
	QueryRing *QueryRing
	Graph     *graph.Graph
	Indexer   *indexer.Indexer
	Config    *config.Config
}

// NewIssueWorker creates a new issue worker
func NewIssueWorker(c *config.Config) *IssueWorker {
	return &IssueWorker{
		QueryRing: NewQueryRing(c.Issue.QueryRing),
		Graph:     graph.New(c),
		Indexer:   indexer.New(c),
		Config:    c,
	}
}

// RecurseIndex with error intervals
func (w *IssueWorker) RecurseIndex() error {
	err := w.Index()

	if err != nil {
		glog.Info("Rolling back query ring")
		w.QueryRing.Rollback()
		w.wait()
	}

	return w.RecurseIndex()
}

// Index GraphQL nodes into Elastic documents
func (w *IssueWorker) Index() error {
	mapping, err := w.Indexer.GetMapping("issue.json")

	if err != nil {
		glog.Errorf("Failed loading issue mapping: %v", err)
		return err
	}

	err = w.Indexer.Ensure(w.Config.Issue.Index, mapping)
	if err != nil {
		glog.Errorf("Failed ensuring index exists: %v", err)
		return err
	}

	return w.indexNextQuery()
}

// Delete index from Elastic cluster
func (w *IssueWorker) Delete() error {
	err := w.Indexer.Delete(w.Config.Issue.Index)
	if err != nil {
		glog.Errorf("Failed deleting index: %v", err)
		return err
	}

	return nil
}

func (w *IssueWorker) indexQuery(query string, endCursor *string, indexCount, parseCount int) error {
	issues, err := w.Graph.FetchIssues(query, endCursor)

	if err != nil {
		glog.Errorf("Failed fetching issues: %v", err)
		return err
	}

	for _, edge := range issues.Data.Search.Edges {
		if len(edge.Node.ID) == 0 {
			continue
		}

		indexed, err := w.reIndexNode(edge.Node)

		if err != nil {
			glog.Errorf("Failed indexing issue: %v", err)
			return err
		}

		parseCount++
		if indexed {
			indexCount++
		}
	}

	if len(issues.Data.Search.PageInfo.EndCursor) > 0 {
		return w.indexQuery(query, &issues.Data.Search.PageInfo.EndCursor, indexCount, parseCount)
	}

	glog.Infof("Indexed, parsed, fetched: %d, %d, %d", indexCount, parseCount, issues.Data.Search.IssueCount)

	if w.Config.Issue.Interval >= 0 {
		w.wait()
		return w.indexNextQuery()
	}

	return nil
}

func (w *IssueWorker) reIndexNode(node gql.Issue) (bool, error) {
	b, err := w.Indexer.Get(
		w.Config.Issue.Index,
		issueType,
		node.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), elasticErrorNotFound) {
			return w.indexDoc(w.parseIssue(node))
		}
		return false, err
	}

	var current mapping.Issue
	err = json.Unmarshal(b, &current)

	if err != nil {
		glog.Errorf("Failed parsing document: %v", err)
		return false, err
	}

	doc := w.parseIssue(node)
	if !cmp.Equal(&current, doc) {
		return w.indexDoc(doc)
	}

	return false, nil
}

func (w *IssueWorker) indexDoc(doc *mapping.Issue) (bool, error) {
	err := w.Indexer.Index(
		w.Config.Issue.Index,
		issueType,
		doc.ID,
		doc,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (w *IssueWorker) indexNextQuery() error {
	query := w.QueryRing.Next()
	glog.Infof("Next query: %s", query)
	return w.indexQuery(query, nil, 0, 0)
}

func (w *IssueWorker) wait() {
	time.Sleep(time.Duration(w.Config.Issue.Interval) * time.Second)
}

func (w *IssueWorker) getReaction(key string, groups []gql.ReactionGroup) int {
	for _, g := range groups {
		if key == g.Content {
			return g.Users.TotalCount
		}
	}
	return 0
}

func (w *IssueWorker) getRepoOwnerLogin(repoURL string) string {
	s := strings.Index(repoURL, "m/")
	e := strings.LastIndex(repoURL, "/")
	if s == -1 || e == -1 {
		return missingValue
	}
	return repoURL[s+2 : e]
}

func (w *IssueWorker) getRepoName(repoURL string) string {
	s := strings.LastIndex(repoURL, "/")
	if s == -1 {
		return missingValue
	}
	return repoURL[s+1:]
}

func (w *IssueWorker) getValue(s string) string {
	if len(s) == 0 {
		return missingValue
	}
	return s
}

func (w *IssueWorker) parseIssue(node gql.Issue) *mapping.Issue {
	authorLogin := w.getValue(node.Author.Login)
	repoOwnerLogin := w.getRepoOwnerLogin(node.Repository.URL)
	repoName := w.getRepoName(node.Repository.URL)
	repoLanguage := w.getValue(node.Repository.PrimaryLanguage.Name)

	return &mapping.Issue{
		ID:                    node.ID,
		URL:                   node.URL,
		Number:                node.Number,
		Title:                 node.Title,
		BodyText:              node.BodyText,
		State:                 node.State,
		ThumbsUp:              w.getReaction(reactionThumbsUp, node.ReactionGroups),
		ThumbsDown:            w.getReaction(reactionThumbsDown, node.ReactionGroups),
		Laugh:                 w.getReaction(reactionLaugh, node.ReactionGroups),
		Hooray:                w.getReaction(reactionHooray, node.ReactionGroups),
		Confused:              w.getReaction(reactionConfused, node.ReactionGroups),
		Heart:                 w.getReaction(reactionHeart, node.ReactionGroups),
		AuthorLogin:           authorLogin,
		AuthorLoginSuggest:    authorLogin,
		RepoOwnerLogin:        repoOwnerLogin,
		RepoOwnerLoginSuggest: repoOwnerLogin,
		RepoName:              repoName,
		RepoNameSuggest:       repoName,
		RepoLanguage:          repoLanguage,
		RepoLanguageSuggest:   repoLanguage,
		RepoForks:             node.Repository.Forks.TotalCount,
		RepoStargazers:        node.Repository.Stargazers.TotalCount,
		CreatedAt:             node.CreatedAt,
		UpdatedAt:             node.UpdatedAt,
	}
}

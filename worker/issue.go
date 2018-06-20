package worker

import (
	"strings"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/graph"
	"github.com/curated/octograph/indexer"
	"github.com/golang/glog"
)

const (
	reactionThumbsUp   = "THUMBS_UP"
	reactionThumbsDown = "THUMBS_DOWN"
	reactionLaugh      = "LAUGH"
	reactionHooray     = "HOORAY"
	reactionConfused   = "CONFUSED"
	reactionHeart      = "HEART"
)

// IssueWorker struct
type IssueWorker struct {
	Graph   *graph.Graph
	Indexer *indexer.Indexer
}

// NewIssueWorker creates a new issue worker
func NewIssueWorker(c *config.Config) *IssueWorker {
	return &IssueWorker{
		Graph:   graph.New(c),
		Indexer: indexer.New(c),
	}
}

// Process GraphQL nodes into Elastic documents
func (w *IssueWorker) Process(query string) error {
	mapping, err := w.Indexer.IssueMapping()

	if err != nil {
		glog.Errorf("Failed loading issue mapping: %v", err)
		return err
	}

	err = w.Indexer.Ensure(indexer.IssueIndex, mapping)
	if err != nil {
		glog.Errorf("Failed ensuring index exists: %v", err)
		return err
	}

	return w.processCursor(query, nil)
}

func (w *IssueWorker) processCursor(query string, endCursor *string) error {
	glog.Infof("Processing cursor: %v", endCursor)

	graphIssues, err := w.Graph.FetchIssues(query, endCursor)
	if err != nil {
		glog.Errorf("Failed fetching issues: %v", err)
		return err
	}

	for _, edge := range graphIssues.Data.Search.Edges {
		if len(edge.Node.ID) == 0 {
			continue
		}

		doc := w.parseIssue(edge.Node)
		err = w.Indexer.Index(
			indexer.IssueIndex,
			indexer.IssueType,
			doc.ID,
			doc,
		)

		if err != nil {
			glog.Errorf("Failed indexing issue: %v", err)
			return err
		}
	}

	if len(graphIssues.Data.Search.PageInfo.EndCursor) > 0 {
		return w.processCursor(query, &graphIssues.Data.Search.PageInfo.EndCursor)
	}

	return nil
}

func (w *IssueWorker) getReaction(key string, groups []graph.ReactionGroup) int {
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
		return ""
	}
	return repoURL[s+2 : e]
}

func (w *IssueWorker) getRepoName(repoURL string) string {
	s := strings.LastIndex(repoURL, "/")
	if s == -1 {
		return ""
	}
	return repoURL[s+1:]
}

func (w *IssueWorker) parseIssue(node graph.Issue) *indexer.Issue {
	return &indexer.Issue{
		ID:             node.ID,
		URL:            node.URL,
		Number:         node.Number,
		Title:          node.Title,
		BodyText:       node.BodyText,
		State:          node.State,
		ThumbsUp:       w.getReaction(reactionThumbsUp, node.ReactionGroups),
		ThumbsDown:     w.getReaction(reactionThumbsDown, node.ReactionGroups),
		Laugh:          w.getReaction(reactionLaugh, node.ReactionGroups),
		Hooray:         w.getReaction(reactionHooray, node.ReactionGroups),
		Confused:       w.getReaction(reactionConfused, node.ReactionGroups),
		Heart:          w.getReaction(reactionHeart, node.ReactionGroups),
		AuthorLogin:    node.Author.Login,
		RepoOwnerLogin: w.getRepoOwnerLogin(node.Repository.URL),
		RepoName:       w.getRepoName(node.Repository.URL),
		RepoLanguage:   node.Repository.PrimaryLanguage.Name,
		RepoForks:      node.Repository.Forks.TotalCount,
		RepoStargazers: node.Repository.Stargazers.TotalCount,
		CreatedAt:      node.CreatedAt,
		UpdatedAt:      node.UpdatedAt,
	}
}

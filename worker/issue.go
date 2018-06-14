package worker

import (
	"log"

	"github.com/curated/octograph/graph"
	"github.com/curated/octograph/indexer"
	"github.com/curated/octograph/logger"
)

// NewIssueWorker creates a new issue worker
func NewIssueWorker() *IssueWorker {
	return &IssueWorker{
		Graph:   graph.New(),
		Indexer: indexer.New(),
		Logger:  logger.New(),
	}
}

// IssueWorker struct
type IssueWorker struct {
	Graph   *graph.Graph
	Indexer *indexer.Indexer
	Logger  *log.Logger
}

// Process GraphQL nodes into Elastic documents
func (w *IssueWorker) Process() error {
	mapping, err := indexer.IssueMapping()
	if err != nil {
		return err
	}
	err = w.Indexer.Ensure(indexer.IssueIndex, mapping)
	if err != nil {
		return err
	}
	return w.processCursor(nil)
}

func (w *IssueWorker) processCursor(endCursor *string) error {
	w.Logger.Printf("Processing cursor: %v", endCursor)
	graphIssues, err := w.Graph.FetchIssues(endCursor)
	if err != nil {
		return err
	}
	for _, edge := range graphIssues.Data.Search.Edges {
		if len(edge.Node.ID) == 0 {
			continue
		}
		doc := w.parse(edge.Node)
		err = w.Indexer.Index(
			indexer.IssueIndex,
			indexer.IssueType,
			doc.ID,
			doc,
		)
		if err != nil {
			return err
		}
	}
	if len(graphIssues.Data.Search.PageInfo.EndCursor) > 0 {
		return w.processCursor(&graphIssues.Data.Search.PageInfo.EndCursor)
	}
	return nil
}

func (w *IssueWorker) parse(node graph.Issue) *indexer.Issue {
	return &indexer.Issue{
		ID:              node.ID,
		URL:             node.URL,
		Number:          node.Number,
		Title:           node.Title,
		BodyText:        node.BodyText,
		State:           node.State,
		ThumbsUp:        w.get("THUMBS_UP", node.ReactionGroups),
		ThumbsDown:      w.get("THUMBS_DOWN", node.ReactionGroups),
		Laugh:           w.get("LAUGH", node.ReactionGroups),
		Hooray:          w.get("HOORAY", node.ReactionGroups),
		Confused:        w.get("CONFUSED", node.ReactionGroups),
		Heart:           w.get("HEART", node.ReactionGroups),
		AuthorID:        node.Author.ID,
		AuthorURL:       node.Author.URL,
		AuthorLogin:     node.Author.Login,
		AuthorAvatar:    node.Author.AvatarURL,
		RepoID:          node.Repository.ID,
		RepoURL:         node.Repository.URL,
		RepoName:        node.Repository.Name,
		RepoLanguage:    node.Repository.PrimaryLanguage.Name,
		RepoForks:       node.Repository.Forks.TotalCount,
		RepoStargazers:  node.Repository.Stargazers.TotalCount,
		RepoOwnerID:     node.Repository.Owner.ID,
		RepoOwnerURL:    node.Repository.Owner.URL,
		RepoOwnerLogin:  node.Repository.Owner.Login,
		RepoOwnerAvatar: node.Repository.Owner.AvatarURL,
		CreatedAt:       node.CreatedAt,
		UpdatedAt:       node.UpdatedAt,
	}
}

func (w *IssueWorker) get(key string, groups []graph.ReactionGroup) int {
	for _, g := range groups {
		if key == g.Content {
			return g.Users.TotalCount
		}
	}
	return 0
}

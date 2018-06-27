package parser

import (
	"strings"

	"github.com/curated/octograph/gql"
	"github.com/curated/octograph/mapping"
)

const (
	reactionThumbsUp   = "THUMBS_UP"
	reactionThumbsDown = "THUMBS_DOWN"
	reactionLaugh      = "LAUGH"
	reactionHooray     = "HOORAY"
	reactionConfused   = "CONFUSED"
	reactionHeart      = "HEART"

	ownerSep   = "m/"
	repoSep    = "/"
	missingVal = "?"
)

// ParseIssue from GraphQL node into Elastic document
func ParseIssue(node gql.Issue) *mapping.Issue {
	return &mapping.Issue{
		Number:         node.Number,
		Title:          node.Title,
		BodyText:       node.BodyText,
		State:          node.State,
		ThumbsUp:       getReaction(reactionThumbsUp, node.ReactionGroups),
		ThumbsDown:     getReaction(reactionThumbsDown, node.ReactionGroups),
		Laugh:          getReaction(reactionLaugh, node.ReactionGroups),
		Hooray:         getReaction(reactionHooray, node.ReactionGroups),
		Confused:       getReaction(reactionConfused, node.ReactionGroups),
		Heart:          getReaction(reactionHeart, node.ReactionGroups),
		AuthorLogin:    getValue(node.Author.Login),
		RepoOwnerLogin: getRepoOwnerLogin(node.Repository.URL),
		RepoName:       getRepoName(node.Repository.URL),
		RepoLanguage:   getValue(node.Repository.PrimaryLanguage.Name),
		RepoForks:      node.Repository.Forks.TotalCount,
		RepoStargazers: node.Repository.Stargazers.TotalCount,
		CreatedAt:      node.CreatedAt,
		UpdatedAt:      node.UpdatedAt,
	}
}

func getReaction(key string, groups []gql.ReactionGroup) int {
	for _, g := range groups {
		if key == g.Content {
			return g.Users.TotalCount
		}
	}
	return 0
}

func getRepoOwnerLogin(repoURL string) string {
	s := strings.Index(repoURL, ownerSep)
	e := strings.LastIndex(repoURL, repoSep)
	if s == -1 || e == -1 {
		return missingVal
	}
	return repoURL[s+2 : e]
}

func getRepoName(repoURL string) string {
	s := strings.LastIndex(repoURL, repoSep)
	if s == -1 {
		return missingVal
	}
	return repoURL[s+1:]
}

func getValue(s string) string {
	if len(s) == 0 {
		return missingVal
	}
	return s
}

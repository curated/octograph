package indexer_test

import (
	"encoding/json"
	"testing"

	"github.com/curated/octograph/indexer"

	"github.com/stretchr/testify/assert"
)

func TestIssueMapping(t *testing.T) {
	issueMapping, err := indexer.IssueMapping()
	assert.Nil(t, err)

	var jsonMapping map[string]interface{}
	err = json.Unmarshal([]byte(issueMapping), &jsonMapping)
	assert.Nil(t, err)

	assert.True(t, len(issueMapping) > 0)
}

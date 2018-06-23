package config_test

import (
	"os"
	"testing"

	"github.com/curated/octograph/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := config.New("../")

	assert.Equal(t, "test", c.Env)
	assert.Equal(t, "../", c.Root)
	assert.NotEmpty(t, c.Elastic.URL)
	assert.NotEmpty(t, c.Elastic.Username)
	assert.NotEmpty(t, c.Elastic.Password)
	assert.NotEmpty(t, c.GitHub.Token)
	assert.Equal(t, "reactions:>=3000", c.IssueWorker.Query)
}

func TestOverride(t *testing.T) {
	orig := os.Getenv("CONFIG")
	os.Setenv("CONFIG", "config/config.sample.json")
	c := config.New("../")

	assert.Equal(t, "sample", c.Env)
	assert.Equal(t, "../", c.Root)
	assert.Empty(t, c.Elastic.URL)
	assert.Empty(t, c.Elastic.Username)
	assert.Empty(t, c.Elastic.Password)
	assert.Empty(t, c.GitHub.Token)
	assert.Equal(t, "reactions:>=100", c.IssueWorker.Query)

	os.Setenv("CONFIG", orig)
}

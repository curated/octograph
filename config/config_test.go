package config_test

import (
	"strings"
	"testing"

	"github.com/curated/octograph/config"
	"github.com/stretchr/testify/assert"
)

func TestConstructor(t *testing.T) {
	c := config.New()
	assert.True(t, len(c.Elastic.URL) > 0)
	assert.True(t, len(c.Elastic.Username) > 0)
	assert.True(t, len(c.Elastic.Password) > 0)
	assert.True(t, len(c.GitHub.Token) > 0)
}

func TestGetPath(t *testing.T) {
	main := config.GetPath("main.go")
	ct := config.GetPath("config/config_test.go")
	assert.True(t, strings.Contains(main, "/src/github.com/curated/octograph/main.go"))
	assert.True(t, strings.Contains(ct, "/src/github.com/curated/octograph/config/config_test.go"))
}

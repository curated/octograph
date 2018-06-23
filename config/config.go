package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/glog"
)

const (
	configKey   = "CONFIG"
	defaultFile = "config/dev.config.json"
)

// Config values
type Config struct {
	Root string

	Elastic struct {
		URL      string
		Username string
		Password string
	}

	GitHub struct {
		Token string
	}

	IssueWorker struct {
		Query string
	}
}

// New creates config from ENV or default file at relative path to package root
func New(root string) *Config {
	c := Config{
		Root: root,
	}

	f := os.Getenv(configKey)
	if len(f) == 0 {
		f = defaultFile
	}

	b, err := ioutil.ReadFile(c.GetPath(f))
	if err != nil {
		glog.Fatalf("Failed loading '%s' with error: %v", f, err)
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		glog.Fatalf("Failed parsing '%s' with error: %v", f, err)
	}

	return &c
}

// GetPath returns relative path within package root
func (c *Config) GetPath(path string) string {
	return filepath.Join(c.Root, path)
}

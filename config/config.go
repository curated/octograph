package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config values
type Config struct {
	Elastic struct {
		URL      string
		Username string
		Password string
	}
	GitHub struct {
		Token string
	}
}

// Load config from ENV or default file
func Load() *Config {
	f := filename()
	return parse(read(f), f)
}

func filename() string {
	f := os.Getenv("CONFIG")
	if len(f) == 0 {
		return "./config/config.json"
	}
	return f
}

func read(f string) []byte {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Printf("Failed loadig '%s' with error: %v", f, err)
		os.Exit(1)
	}
	return b
}

func parse(b []byte, f string) *Config {
	var c Config
	err := json.Unmarshal(b, &c)
	if err != nil {
		fmt.Printf("Failed parsing '%s' with error: %v", f, err)
		os.Exit(1)
	}
	return &c
}

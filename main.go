package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/worker"
)

const root = "./"
const query = "reactions:>=100"

func main() {
	flag.Parse()
	c := config.New(root)
	err := worker.NewIssueWorker(c).Process(query)

	if err != nil {
		glog.Fatalf("Failed running issue worker: %v", err)
	}
}

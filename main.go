package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/curated/octograph/worker"
)

func main() {
	flag.Parse()
	err := worker.NewIssueWorker().Process("reactions:>=100")
	if err != nil {
		glog.Fatalf("Failed running issue worker: %v", err)
	}
}

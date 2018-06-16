package main

import (
	"flag"
	"os"

	"github.com/curated/octograph/worker"
)

func main() {
	flag.Parse()
	err := worker.NewIssueWorker().Process()
	if err != nil {
		os.Exit(1)
	}
}

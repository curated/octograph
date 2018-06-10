package main

import (
	"os"

	"github.com/curated/octograph/worker"
)

func main() {
	err := worker.NewIssueWorker().Process()
	if err != nil {
		os.Exit(1)
	}
}

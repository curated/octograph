package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/curated/octograph/config"
	"github.com/curated/octograph/worker"
)

const (
	root          = "./"
	processFlag   = "process"
	processUsage  = "Processing method: [index|delete] defaults to index"
	indexProcess  = "index"
	deleteProcess = "delete"
	query         = "reactions:>=25"
)

func main() {
	var err error
	c := config.New(root)
	issueWorker := worker.NewIssueWorker(c)
	process := flag.String(processFlag, indexProcess, processUsage)
	flag.Parse()

	switch *process {
	case indexProcess:
		err = issueWorker.Index(query)
	case deleteProcess:
		err = issueWorker.Delete()
	}

	if err != nil {
		glog.Fatalf("Failed running issue worker: %v", err)
	}
}

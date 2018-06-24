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
)

func main() {
	var err error
	c := config.New(root)
	issueWorker := worker.NewIssueWorker(c)
	process := flag.String(processFlag, indexProcess, processUsage)
	flag.Parse()

	glog.Infof("Running %s on %s", *process, c.Env)

	switch *process {
	case indexProcess:
		err = issueWorker.RecurseIndex()
	case deleteProcess:
		err = issueWorker.Delete()
	default:
		glog.Fatalf("Invalid process argument: %s", *process)
	}

	if err != nil {
		glog.Fatalf("Failed running issue worker: %v", err)
	}
}

package logger

import (
	"log"
	"os"
)

// New creates a logger
func New() *log.Logger {
	return log.New(
		os.Stdout,
		"",
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile,
	)
}

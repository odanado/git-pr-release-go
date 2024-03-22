package main

import (
	"log"
	"os"
)

var logger *log.Logger

func GetLogger() *log.Logger {
	if logger != nil {
		return logger
	}

	logger = log.New(os.Stdout, "git-pr-release-go", 0)
	return logger
}

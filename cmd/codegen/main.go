package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func initLog() *log.Logger {
	l := log.New()
	l.SetFormatter(&log.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(log.InfoLevel)

	return l
}

func main() {
	initLog()
}

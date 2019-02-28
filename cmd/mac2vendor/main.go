package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true

	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(customFormatter)
}

func main() {
	http.HandleFunc("/", logger(lookup))
	log.Info("fired up service...")
	http.ListenAndServe(":9000", nil)
}

package main

import (
	"flag"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	m2v "github.com/n3integration/mac2vendor"
)

var mac2vnd *m2v.Mac2Vendor

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true

	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(customFormatter)
}

func main() {
	dat := flag.String("file", m2v.Dat, "mac address to vendor mapping file")
	flag.Parse()

	var err error
	mac2vnd, err = m2v.Load(*dat)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	http.HandleFunc("/", logger(lookup))
	log.Info("fired up service...")
	http.ListenAndServe(":9000", nil)
}

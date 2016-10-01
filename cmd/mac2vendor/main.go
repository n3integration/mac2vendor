package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	m2v "github.com/n3integration/mac2vendor"
	"net/http"
	"os"
)

var mac2vnd *m2v.Mac2Vendor

type Mac2Vnd struct {
	Mac    string `json:"mac,omitempty"`
	Vendor string `json:"vendor,omitempty"`
	Error  string `json:"error,omitempty"`
}

type lookupHandler struct {
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var err error
	mac2vnd, err = m2v.Load(m2v.MAC2VND)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	http.HandleFunc("/", logger(lookup))
	log.Info("fired up...")
	http.ListenAndServe(":9000", nil)
}

func newMac2Vnd(mac, vendor string, err error) *Mac2Vnd {
	if err != nil {
		return &Mac2Vnd{Error: err.Error()}
	}
	return &Mac2Vnd{Mac: mac, Vendor: vendor}
}

func logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RemoteAddr, " ", r.Method, " ", r.RequestURI, " ", r.Proto, " ", r.ContentLength)
		next.ServeHTTP(w, r)
	})
}

func lookup(w http.ResponseWriter, r *http.Request) {
	mac := r.URL.Path[1:]
	vendor, err := mac2vnd.Lookup(mac)
	response := newMac2Vnd(mac, vendor, err)
	json, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Error != "" {
		http.Error(w, string(json), http.StatusNotFound)
		return
	}

	w.Write(json)
}

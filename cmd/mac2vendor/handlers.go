package main

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Mac2Vnd is a resource model
type Mac2Vnd struct {
	Mac    string `json:"mac,omitempty"`
	Vendor string `json:"vendor,omitempty"`
	Error  string `json:"error,omitempty"`
}

// newMac2Vnd initializes a new response
func newMac2Vnd(mac, vendor string, err error) *Mac2Vnd {
	if err != nil {
		return &Mac2Vnd{Error: err.Error()}
	}
	return &Mac2Vnd{Mac: mac, Vendor: vendor}
}

// logger is logging middleware
func logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Info(r.RemoteAddr, " ", r.Method, " ", r.RequestURI, " ", r.Proto, " ", r.ContentLength)
	})
}

// lookup provides the mac address to vendor lookup service handler
func lookup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

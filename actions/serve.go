package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/n3integration/mac2vendor"
	"gopkg.in/urfave/cli.v1"
)

var port uint

func init() {
	log.SetFlags(log.LstdFlags)
	register(cli.Command{
		Name:   "serve",
		Action: serveAction,
		Usage:  "expose mac2vendor as a web service",
		Flags: []cli.Flag{
			cli.UintFlag{
				Name:        "port",
				EnvVar:      "PORT",
				Value:       9000,
				Destination: &port,
				Usage:       "the port to which the service should bind",
			},
		},
	})
}

func serveAction(_ *cli.Context) error {
	http.HandleFunc("/", logger(lookup))
	log.Printf("Service listening at 127.0.0.1:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// Mac2Vnd is a resource model
type Mac2Vnd struct {
	Mac    string `json:"mac,omitempty"`
	Vendor string `json:"vendor,omitempty"`
	Error  error  `json:"error,omitempty"`
}

// newMac2Vnd initializes a new response
func newMac2Vnd(mac, vendor string, err error) *Mac2Vnd {
	if err != nil {
		return &Mac2Vnd{
			Error: err,
		}
	}
	return &Mac2Vnd{
		Mac:    mac,
		Vendor: vendor,
	}
}

type interceptor struct {
	Status   int
	Bytes    int64
	delegate http.ResponseWriter
}

func (i *interceptor) Header() http.Header {
	return i.delegate.Header()
}

func (i *interceptor) WriteHeader(statusCode int) {
	i.Status = statusCode
	i.delegate.WriteHeader(statusCode)
}

func (i *interceptor) Write(p []byte) (n int, err error) {
	i.Bytes += int64(len(p))
	return i.delegate.Write(p)
}

// logger is logging middleware
func logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wi := &interceptor{delegate: w}
		defer func() {
			log.Println(r.RemoteAddr, " ", r.Method, " ", r.RequestURI, " ", r.Proto, " ", wi.Status, " ", wi.Bytes)
		}()
		next.ServeHTTP(wi, r)
	})
}

// lookup provides the mac address to vendor lookup service handler
func lookup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mac := r.URL.Path[1:]
	vendor, err := mac2vendor.Lookup(mac)
	response := newMac2Vnd(mac, vendor, err)
	json, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Error != nil {
		status := http.StatusNotFound
		switch response.Error.(type) {
		case *net.AddrError:
			status = http.StatusBadRequest
		}
		http.Error(w, string(json), status)
		return
	}

	if _, err := w.Write(json); err != nil {
		log.Println("failed to write response: ", err)
	}
}

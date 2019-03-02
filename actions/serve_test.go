package actions

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestServe(t *testing.T) {
	validMAC := "/84:38:35:77:aa:52"
	type args struct {
		r    *http.Request
		code int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Default",
			args: args{
				httptest.NewRequest(http.MethodGet, validMAC, nil),
				http.StatusOK,
			},
		}, {
			name: "No MAC Address Provided",
			args: args{
				httptest.NewRequest(http.MethodGet, "/", nil),
				http.StatusBadRequest,
			},
		}, {
			name: "Unsupported Method",
			args: args{
				httptest.NewRequest(http.MethodPost, validMAC, nil),
				http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			lookup(w, tt.args.r)
			if w.Code != tt.args.code {
				t.Errorf("received unexpected status code: %v; expected %v", w.Code, tt.args.code)
			}
		})
	}
}

func TestLogger(t *testing.T) {
	next := http.NotFound
	out := new(bytes.Buffer)

	log.SetOutput(out)
	logger(next).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if len(out.Bytes()) == 0 {
		t.Errorf("expected request to be logged, output (%v)", len(out.Bytes()))
	} else if !strings.Contains(out.String(), strconv.Itoa(http.StatusNotFound)) {
		t.Errorf("expected logged request to include response status code, but not found: %s", out)
	}
}

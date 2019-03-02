package actions

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpdate(t *testing.T) {
	tplPath = "../templates/mac2vnd.tpl"
	goldenFile := "testdata/oui.golden"
	oui, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal("failed to load golden file: ", err)
	}

	t.Run("DownloadMacTable", func(t *testing.T) {
		dst := "oui.txt"
		defer os.Remove(dst)

		t.Run("Not Found", func(t *testing.T) {
			server := httptest.NewServer(http.NotFoundHandler())
			defer server.Close()

			source = server.URL
			if err := downloadMacTable(dst); err == nil {
				t.Error("Expected download to fail, but no errors were returned")
			}
		})

		t.Run("Downloaded", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				io.Copy(w, bytes.NewReader(oui))
			}))
			defer server.Close()

			source = server.URL
			http.DefaultClient = server.Client()
			if err := downloadMacTable(dst); err != nil {
				t.Error("failed to download mac table: ", err)
			}
		})
	})

	t.Run("Generate Mapping", func(t *testing.T) {
		key := "3c:d9:2b"
		outfile := "mapping.go"
		mapping := map[string]string{
			key: "Hewlett Packard",
		}

		defer os.Remove(outfile)
		if err := generateMapping(mapping); err != nil {
			t.Error("failed to generate mapping file: ", err)
		}

		b, _ := ioutil.ReadFile(outfile)
		if !bytes.Contains(b, []byte(key)) {
			t.Error("failed to find key in mapping file: ", string(b))
		}
	})

	t.Run("Transform", func(t *testing.T) {
		defer os.Remove("mapping.go")
		if err := transform(goldenFile, "oui.txt"); err != nil {
			t.Error("failed to transform oui file: ", err)
		}
	})
}

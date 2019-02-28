package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	// Dat is the default mapping file name
	Dat          = "mac2vnd.dat"
	delimiter    = "\t"
	prefixLength = 5
)

func init() {
	log.SetFlags(log.LstdFlags)
}

// Load initializes the mac to vendor mapping from the provided src file
func load(src string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		log.Println("loading mac2vendor data. please be patient...")
		oui := "/tmp/oui.txt"

		if err := downloadMacTable(oui); err != nil {
			return err
		}

		if err := transform(oui, src); err != nil {
			return err
		}
	}

	return nil
}

// downloadMacTable downloads the oui.txt file from ieee
func downloadMacTable(dst string) error {
	src := "http://standards-oui.ieee.org/oui/oui.txt"

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		log.Println("saving file to", dst)

		response, err := http.Get(src)
		if err != nil {
			return errors.Wrap(err, "failed to download "+src)
		}
		defer response.Body.Close()

		output, err := os.Create(dst)
		if err != nil {
			return errors.Wrap(err, "failed to create "+dst)
		}
		defer output.Close()

		writer := bufio.NewWriter(output)
		defer writer.Flush()

		if _, err = io.Copy(writer, response.Body); err != nil {
			return errors.Wrap(err, "failed to write "+src)
		}
	}
	return nil
}

// transform converts the raw contents of src into key=value form in dst
func transform(src string, dst string) error {
	_, err := os.Stat(dst)
	if os.IsNotExist(err) {
		log.Println("transforming", src, "into", dst)

		f, err := os.Open(src)
		if err != nil {
			return errors.Wrap(err, "failed to parse")
		}
		defer f.Close()

		output, err := os.Create(dst)
		if err != nil {
			return errors.Wrap(err, "failed to parse")
		}
		defer output.Close()

		pattern := regexp.MustCompile(`\s*([0-9a-zA-Z]+)[\s]*\(base 16\)[\s]*([^\r]+)`)
		reader := bufio.NewReader(f)

		writer := bufio.NewWriter(output)
		defer writer.Flush()

		mapping := make(map[string]string)

		for {
			line, err := reader.ReadBytes('\n')

			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			if pattern.Match(line) {
				parts := pattern.FindStringSubmatch(string(line))
				mapping[delimit(strings.ToLower(parts[1]))] = parts[2]
				writer.WriteString(fmt.Sprintf("%s%s%s\n", strings.ToLower(delimit(parts[1])), delimiter, parts[2]))
			}
		}

		defer os.Remove(dst)
		goTemplate, err := ioutil.ReadFile("templates/mac2vnd.tpl")
		if err != nil {
			return errors.Wrap(err, "failed to read template file")
		}
		log.Println("Parsing template...")
		t, err := template.New("t").Parse(string(goTemplate))
		if err != nil {
			return errors.Wrap(err, "failed to parse template")
		}

		log.Println("Executing template...")
		buffer := new(bytes.Buffer)
		err = t.Execute(buffer, mapping)
		if err != nil {
			return errors.Wrap(err, "failed to execute template")
		}

		formatted, err := format.Source(buffer.Bytes())
		if err != nil {
			return errors.Wrap(err, "failed to format source")
		}

		ioutil.WriteFile("mapping.go", formatted, 0755)
	}
	return nil
}

func delimit(prefix string) string {
	var mac bytes.Buffer
	for i, c := range prefix {
		mac.WriteRune(c)
		if i%2 != 0 && i < prefixLength {
			mac.WriteString(":")
		}
	}
	return mac.String()
}

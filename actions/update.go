package actions

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

const (
	datFile      = "mac2vnd.dat"
	delimiter    = "\t"
	prefixLength = 5
)

var (
	source  = "http://standards-oui.ieee.org/oui/oui.txt"
	tplPath = "templates/mac2vnd.tpl"
)

func init() {
	register(cli.Command{
		Name:    "update",
		Action:  updateAction,
		Aliases: []string{"up"},
		Usage:   "update the mac address vendor mapping to the latest oui listing",
	})
}

func updateAction(_ *cli.Context) error {
	return load(datFile)
}

// Load initializes the mac to vendor mapping from the provided src file
func load(src string) error {
	oui := "oui.txt"
	if _, err := os.Stat(src); os.IsNotExist(err) {
		log.Println("loading mac2vendor data...")
		defer os.Remove(oui)

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
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		log.Println("saving file to", dst)

		response, err := http.Get(source)
		if err != nil {
			return errors.Wrap(err, "failed to download "+source)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return errors.New("failed to download " + source + "; server responded with " + response.Status)
		}

		output, err := os.Create(dst)
		if err != nil {
			return errors.Wrap(err, "failed to create "+dst)
		}
		defer output.Close()

		writer := bufio.NewWriter(output)
		defer writer.Flush()

		if _, err = io.Copy(writer, response.Body); err != nil {
			return errors.Wrap(err, "failed to write "+source)
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

		pattern := regexp.MustCompile(`\s*([0-9a-zA-Z]+)[\s]*\(base 16\)[\s]*([^\r\n]+)`)
		reader := bufio.NewReader(f)

		writer := bufio.NewWriter(output)
		defer writer.Flush()

		mapping := make(map[string]string)
		defer os.Remove(dst)

		for {
			line, err := reader.ReadBytes('\n')

			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			if pattern.Match(line) {
				parts := pattern.FindStringSubmatch(string(line))
				prefix := delimit(strings.ToLower(parts[1]))
				mapping[prefix] = parts[2]
				writer.WriteString(fmt.Sprintf("%s%s%s\n", prefix, delimiter, parts[2]))
			}
		}

		if err := generateMapping(mapping); err != nil {
			return err
		}
	}
	return nil
}

func generateMapping(mapping map[string]string) error {
	goTemplate, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return errors.Wrap(err, "failed to read template file")
	}

	log.Println("parsing template...")
	t, err := template.New("t").Parse(string(goTemplate))
	if err != nil {
		return errors.Wrap(err, "failed to parse template")
	}

	log.Println("executing template...")
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, mapping)
	if err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	formatted, err := format.Source(buffer.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to format source")
	}

	return ioutil.WriteFile("mapping.go", formatted, 0755)
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

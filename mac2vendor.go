package mac2vendor

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	// Dat is the default mapping file name
	Dat       = "mac2vnd.dat"
	delimiter = "\t"
)

// Mac2Vendor encapsulates the mac address to vendor lookup
type Mac2Vendor struct {
	mapping map[string]string
}

func init() {
	log.SetFlags(log.LstdFlags)
}

// Load initializes the mac to vendor mapping from the provided src file
func Load(src string) (*Mac2Vendor, error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		log.Println("loading. please be patient...")
		oui := "/tmp/oui.txt"

		if err := downloadMacTable(oui); err != nil {
			return nil, err
		}

		if err := transform(oui, src); err != nil {
			return nil, err
		}
	}

	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("load", "-", err)
	}
	defer f.Close()

	return initCache(bufio.NewReader(f))
}

// Lookup resolves the mac address vendor, if found; otherwise an error is returned
func (m *Mac2Vendor) Lookup(mac string) (string, error) {
	normalized := strings.Replace(strings.Replace(mac, ":", "", -1), "-", "", -1)
	if len(normalized) > 6 {
		normalized = normalized[:6]
	}

	if val, ok := m.mapping[normalized]; ok {
		return val, nil
	}

	return "", fmt.Errorf("%s not found", mac)
}

// downloadMacTable downloads the oui.txt file from ieee
func downloadMacTable(dst string) error {
	src := "http://standards.ieee.org/develop/regauth/oui/oui.txt"

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		log.Println("downloading", src, "to", dst)

		output, err := os.Create(dst)
		if err != nil {
			return fmt.Errorf("creating", dst, "-", err)
		}
		defer output.Close()

		response, err := http.Get(src)
		if err != nil {
			return fmt.Errorf("downloading", src, "-", err)
		}
		defer response.Body.Close()

		writer := bufio.NewWriter(output)
		defer writer.Flush()

		_, err = io.Copy(writer, response.Body)

		if err != nil {
			return fmt.Errorf("downloading", src, "-", err)
		}
	}
	return nil
}

// transform converts the raw contents of src into key=value form in dst
func transform(src string, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		log.Println("transforming", src, "into", dst)

		f, err := os.Open(src)
		if err != nil {
			return fmt.Errorf("parsing", "-", err)
		}
		defer f.Close()

		output, err := os.Create(dst)
		if err != nil {
			return fmt.Errorf("parsing", "-", err)
		}
		defer output.Close()

		pattern := regexp.MustCompile(`\s*([0-9a-zA-Z]+)[\s]*\(base 16\)[\s]*(.*)`)
		reader := bufio.NewReader(f)

		writer := bufio.NewWriter(output)
		defer writer.Flush()

		for {
			line, err := reader.ReadBytes('\n')

			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}

			if pattern.Match(line) {
				parts := pattern.FindStringSubmatch(string(line))
				writer.WriteString(fmt.Sprintf("%s%s%s\n", parts[1], delimiter, parts[2]))
			}
		}
	} else {
		return err
	}
}

// initCache initializes a new Mac2Vendor using the provided Reader
func initCache(r *bufio.Reader) (*Mac2Vendor, error) {
	mapping := make(map[string]string)

	for {
		line, err := r.ReadBytes('\n')

		if err == io.EOF {
			return &Mac2Vendor{mapping}, nil
		} else if err != nil {
			return nil, err
		}

		parts := strings.SplitN(string(line), delimiter, 2)
		mapping[parts[0]] = strings.TrimSpace(parts[1])
	}
}

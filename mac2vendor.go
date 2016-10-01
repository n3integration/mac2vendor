package mac2vendor

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const OUI = "oui.txt"
const MAC2VND = "mac2vnd.dat"

type Mac2Vendor struct {
	mapping map[string]string
}

// load the mac to vendor mapping from the provided src file
func Load(src string) (*Mac2Vendor, error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		oui := fmt.Sprintf("/tmp/%s", OUI)

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

	mapping := make(map[string]string)
	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadBytes('\n')

		if err == io.EOF {
			return &Mac2Vendor{mapping}, nil
		}

		parts := strings.SplitN(string(line), "\t", 2)
		mapping[parts[0]] = strings.TrimSpace(parts[1])
	}
}

// resolve the mac address vendor
func (this *Mac2Vendor) Lookup(mac string) (string, error) {
	normalized := strings.Replace(strings.Replace(mac, ":", "", -1), "-", "", -1)
	if len(normalized) > 6 {
		normalized = normalized[:6]
	}

	if val, ok := this.mapping[normalized]; ok {
		return val, nil
	}
	return "", fmt.Errorf("not found")
}

// download the oui.txt file from ieee
func downloadMacTable(dst string) error {
	src := "http://standards.ieee.org/develop/regauth/oui/oui.txt"

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		fmt.Println("downloading", src, "to", dst)

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

// transform the raw contents of src into key=value form in dst
func transform(src string, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		fmt.Println("transforming", src, "into", dst)

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
			}

			if pattern.Match(line) {
				parts := pattern.FindStringSubmatch(string(line))
				writer.WriteString(fmt.Sprintf("%s\t%s\n", parts[1], parts[2]))
			}
		}
	}
	return nil
}

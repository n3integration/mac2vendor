package mac2vendor

import (
	"github.com/pkg/errors"
	"net"
	"strings"
)

var (
	errCannotResolveType = errors.New("cannot resolve type to mac address")
	mapping              = make(map[string]string)
)

// IsLoaded is a predicate to determine whether or not the mapping table was loaded
func IsLoaded() bool {
	return len(mapping) > 0
}

// Lookup resolves the provided MAC address to the registered vendor
func Lookup(v interface{}) (string, error) {
	var mac net.HardwareAddr
	switch v.(type) {
	case string:
		var err error
		mac, err = net.ParseMAC(v.(string))
		if err != nil {
			return "", err
		}
	case net.HardwareAddr:
		mac = v.(net.HardwareAddr)
	default:
		return "", errCannotResolveType
	}

	prefix := mac[:3].String()
	if val, ok := mapping[strings.ToLower(prefix)]; ok {
		return val, nil
	}

	return "", nil
}

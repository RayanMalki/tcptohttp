package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var rn = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Get(key string) string {
	key = strings.ToLower(key)
	if val, ok := h[key]; ok {
		return val
	}
	return ""
}
func IsKeyCharValid(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			strings.ContainsRune("!#$%&'*+-.^_`|~", r) {
			continue
		}
		return false
	}
	return true
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	rawName := parts[0]
	rawValue := parts[1]

	if len(rawName) == 0 || bytes.HasPrefix(rawName, []byte(" ")) || bytes.HasSuffix(rawName, []byte(" ")) || bytes.Contains(rawName, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	name := string(rawName)
	value := strings.TrimSpace(string(rawValue))
	return name, value, nil

}
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	read := 0
	done = false
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break

		}
		//Empty header
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err

		}
		read += idx + len(rn)

		if !IsKeyCharValid(name) {

			return 0, false, fmt.Errorf("key contains invalid character in header key")
		}

		key := strings.ToLower(name)
		if oldValue, exists := h[key]; exists {
			h[key] = oldValue + "," + value
		} else {
			h[key] = value
		}

	}
	return read, done, nil
}

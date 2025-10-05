package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// RequestFromReader reads the full HTTP request from an io.Reader
// and parses only the Request-Line.
func RequestFromReader(reader io.Reader) (*Request, error) {
	// Read everything from the reader into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Convert bytes to string
	raw := string(data)

	// Split the request by the first newline (\r\n)
	lines := strings.SplitN(raw, "\r\n", 2)
	if len(lines) == 0 {
		return nil, errors.New("empty request")
	}

	// Parse only the first line
	requestLineStr := lines[0]
	reqLine, err := parseRequestLine(requestLineStr)
	if err != nil {
		return nil, err
	}

	// Return the parsed request
	return &Request{RequestLine: reqLine}, nil
}

// parseRequestLine parses a single HTTP request line, e.g. "GET / HTTP/1.1"
func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return RequestLine{}, errors.New("invalid number of parts in request line")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	// Validate method: must be all uppercase letters
	for _, r := range method {
		if !unicode.IsUpper(r) {
			return RequestLine{}, errors.New("invalid method format")
		}
	}

	// Validate version: must be "HTTP/1.1"
	if !strings.HasPrefix(version, "HTTP/") {
		return RequestLine{}, errors.New("invalid HTTP version format")
	}
	if strings.TrimPrefix(version, "HTTP/") != "1.1" {
		return RequestLine{}, errors.New("unsupported HTTP version")
	}

	// Construct the RequestLine struct
	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, nil
}

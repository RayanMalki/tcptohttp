package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Request represents a full HTTP request (we're focusing on the request line for now)
type Request struct {
	RequestLine RequestLine
	state       int
}

// RequestLine holds the three components of the HTTP request line
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// Enum (int) for parser state
const (
	stateInitialized = iota
	stateDone
)

// parseRequestLine parses the request line and returns the number of bytes consumed.
// If no full line is found (no \r\n), returns 0 and no error.
func parseRequestLine(data []byte) (RequestLine, int, error) {
	// Look for the end of the request line (CRLF)
	index := bytes.Index(data, []byte("\r\n"))
	if index == -1 {
		// No full line yet
		return RequestLine{}, 0, nil
	}

	line := string(data[:index])
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("invalid request line: must contain method, target, version")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	// Validate method: must be uppercase letters only
	for _, ch := range method {
		if ch < 'A' || ch > 'Z' {
			return RequestLine{}, 0, fmt.Errorf("invalid method: %s", method)
		}
	}

	// Validate version format: HTTP/1.1 only
	if !strings.HasPrefix(version, "HTTP/") {
		return RequestLine{}, 0, fmt.Errorf("invalid version format: %s", version)
	}

	versionNumber := strings.TrimPrefix(version, "HTTP/")
	if versionNumber != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("unsupported HTTP version: %s", versionNumber)
	}

	// Return a parsed request line and how many bytes we consumed (line + CRLF)
	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   versionNumber,
	}, index + 2, nil
}

// parse processes chunks of bytes and updates the request state
func (r *Request) parse(data []byte) (int, error) {
	if r.state == stateDone {
		return 0, nil
	}

	reqLine, consumed, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if consumed == 0 {
		// Not enough data to parse yet
		return 0, nil
	}

	r.RequestLine = reqLine
	r.state = stateDone
	return consumed, nil
}

// RequestFromReader reads from a stream (io.Reader) and builds a Request struct
func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{state: stateInitialized}
	buffer := make([]byte, 0, 8)
	tmp := make([]byte, 8)

	for {
		n, err := reader.Read(tmp)
		if n > 0 {
			buffer = append(buffer, tmp[:n]...)

			consumed, parseErr := r.parse(buffer)
			if parseErr != nil {
				return nil, parseErr
			}

			if r.state == stateDone {
				break
			}

			if consumed > 0 {
				// Remove parsed data from buffer
				buffer = buffer[consumed:]
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	if r.state != stateDone {
		return nil, errors.New("incomplete request")
	}

	return r, nil
}

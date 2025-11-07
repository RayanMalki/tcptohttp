package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/RayanMalki/tcptohttp/internal/headers"
)

// Request represents a full HTTP request (we're focusing on the request line for now)
type Request struct {
	Method      string
	Path        string
	Version     string
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
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
	requestStateStart = iota
	requestStateParsingRequestLine
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
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
	validMethods := []string{"GET", "POST", "HEAD", "PUT", "DELETE", "OPTIONS"}
	isValid := false
	for _, m := range validMethods {
		if method == m {
			isValid = true
			break
		}
	}
	if !isValid {
		return RequestLine{}, 0, fmt.Errorf("invalid HTTP method: %s", method)
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
func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateParsingRequestLine:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return n, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = reqLine
		r.state = requestStateParsingHeaders
		return n, nil

	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return n, err
		}
		if done {
			if r.Headers.Get("content-length") != "" {
				r.state = requestStateParsingBody
			} else {
				r.state = requestStateDone
			}
		}
		return n, nil

	case requestStateParsingBody:
		contentLengthStr := r.Headers.Get("content-length")
		if contentLengthStr == "" {
			r.state = requestStateDone
			return 0, nil
		}

		// Convert Content-Length to int
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, fmt.Errorf("invalid Content-Length: %v", err)
		}

		// Calculate how many bytes are available for body
		remaining := contentLength - len(r.Body)
		if remaining <= 0 {
			r.state = requestStateDone
			return 0, nil
		}

		// Determine how much data we can read from this chunk
		toRead := len(data)
		if toRead > remaining {
			toRead = remaining
		}

		// Append body bytes
		r.Body = append(r.Body, data[:toRead]...)

		// Check if we read the full body
		if len(r.Body) == contentLength {
			r.state = requestStateDone
		}

		return toRead, nil

	default:
		return 0, fmt.Errorf("invalid parser state: %v", r.state)
	}
}

// parse processes chunks of bytes and updates the request state
func (r *Request) Parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])

		if err != nil {
			return totalBytesParsed, err
		}
		if n == 0 {
			break
		}

		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

// RequestFromReader reads from a stream (io.Reader) and builds a Request struct
func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{state: requestStateParsingRequestLine, Headers: headers.NewHeaders()}
	buffer := make([]byte, 0, 8)
	tmp := make([]byte, 8)

	for {
		n, err := reader.Read(tmp)
		if n > 0 {
			buffer = append(buffer, tmp[:n]...)

			consumed, parseErr := r.Parse(buffer)
			if parseErr != nil {
				return nil, parseErr
			}

			if r.state == requestStateDone {
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

	if r.state != requestStateDone {
		return nil, errors.New("incomplete request")
	}

	return r, nil
}

package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

// Good Request line
func TestGoodRequestLine(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

// Good Request line with path
func TestGoodRequestLineWithPath(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

// Good POST Request with path
func TestGoodPOSTRequestWithPath(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("POST /login HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/login", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

// Invalid number of parts in request line
func TestInvalidNumberOfPartsInRequestLine(t *testing.T) {
	_, err := RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.Error(t, err)
}

// Invalid method (out of order) Request line
func TestInvalidMethodRequestLine(t *testing.T) {
	_, err := RequestFromReader(strings.NewReader("GeT /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.Error(t, err)
}

// Invalid version in Request line
func TestInvalidVersionRequestLine(t *testing.T) {
	_, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\n\r\n"))
	require.Error(t, err)
}

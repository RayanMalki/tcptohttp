package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/RayanMalki/tcptohttp/internal/headers"
)

type StatusCode int

const (
	StatusOK            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var reason string

	switch statusCode {
	case StatusOK:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalError:
		reason = "Internal Server Error"
	default:
		reason = ""
	}

	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, reason)

	return err

}

func GetDefaultHeaders(contentLen int) headers.Headers {

	headersMap := headers.NewHeaders()

	lenghthStr := strconv.Itoa(contentLen)

	headersMap["Content-Length"] = lenghthStr
	headersMap["Connection"] = "close"
	headersMap["Content-Type"] = "text/plain"

	return headersMap

}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for key, value := range h {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w, "\r\n")
	return err
}

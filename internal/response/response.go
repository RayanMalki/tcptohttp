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

type Writer struct {
	conn  io.Writer
	state string //"init", "status_written", "headers_written"
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		conn:  conn,
		state: "init",
	}
}

func (w *Writer) WriteStatusLine(code StatusCode) error {
	if w.state != "init" {
		return fmt.Errorf("status already written")
	}

	reason := ""
	switch code {
	case StatusOK:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalError:
		reason = "Internal Server Error"
	default:
		reason = ""
	}

	// Dynamically write the status line
	if _, err := fmt.Fprintf(w.conn, "HTTP/1.1 %d %s\r\n", code, reason); err != nil {
		return err
	}

	w.state = "status_written"
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {

	headersMap := headers.NewHeaders()

	lenghthStr := strconv.Itoa(contentLen)

	headersMap["Content-Length"] = lenghthStr
	headersMap["Connection"] = "close"
	headersMap["Content-Type"] = "text/plain"

	return headersMap

}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != "status_written" {
		return fmt.Errorf("must write status line before headers")
	}

	for key, value := range headers {
		if _, err := fmt.Fprintf(w.conn, "%s: %s\r\n", key, value); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w.conn, "\r\n"); err != nil {
		return err
	}

	w.state = "headers_written"
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != "headers_written" {
		return 0, fmt.Errorf("must write headers before body")
	}
	bytesWritten, err := w.conn.Write(p)

	if err != nil {
		return bytesWritten, err
	}

	w.state = "body_written"

	return bytesWritten, nil
}

package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/RayanMalki/tcptohttp/internal/headers"
	"github.com/RayanMalki/tcptohttp/internal/request"
	"github.com/RayanMalki/tcptohttp/internal/response"
	"github.com/RayanMalki/tcptohttp/internal/server"
)

const port = 42069

func myHandler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video/") {
		filename := strings.TrimPrefix(req.RequestLine.RequestTarget, "/video/")
		filepath := "assets/" + filename

		data, err := os.ReadFile(filepath)
		if err != nil {
			w.WriteStatusLine(response.StatusBadRequest)
			headers := response.GetDefaultHeaders(0)
			headers.Set("Content-Type", "text/html")
			w.WriteHeaders(headers)
			w.WriteBody([]byte("<html><body><h1>Video Not Found</h1></body></html>"))
			return
		}

		w.WriteStatusLine(response.StatusOK)
		headers := response.GetDefaultHeaders(len(data))
		headers.Set("Content-Type", "video/mp4")
		w.WriteHeaders(headers)
		w.WriteBody(data)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
		url := "https://httpbin.org/" + path

		newReq, _ := http.NewRequest("GET", url, nil)
		newReq.Header.Set("Accept-Encoding", "identity")

		client := &http.Client{}
		resp, err := client.Do(newReq)
		if err != nil {
			w.WriteStatusLine(response.StatusInternalError)
			headers := response.GetDefaultHeaders(0)
			headers.Set("Content-Type", "text/html")
			w.WriteHeaders(headers)
			w.WriteBody([]byte("<html><body><h1>Proxy Error</h1></body></html>"))
			return
		}
		defer resp.Body.Close()

		w.WriteStatusLine(response.StatusOK)
		hdrs := response.GetDefaultHeaders(0)
		delete(hdrs, "Content-Length")
		hdrs.Set("Transfer-Encoding", "chunked")
		hdrs.Set("Trailer", "X-Content-SHA256, X-Content-Length")
		hdrs.Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeaders(hdrs)

		var fullBody []byte
		buf := make([]byte, 1024)

		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				w.WriteChunkedBody(buf[:n])
				fullBody = append(fullBody, buf[:n]...)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
		}

		w.WriteChunkedBodyDone()

		hash := sha256.Sum256(fullBody)
		trailers := headers.NewHeaders()
		trailers["X-Content-SHA256"] = fmt.Sprintf("%x", hash)
		trailers["X-Content-Length"] = fmt.Sprintf("%d", len(fullBody))
		w.WriteTrailers(trailers)
		return
	}

	var html string
	var status response.StatusCode

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.StatusBadRequest
		html = `<html><body><h1>Bad Request</h1></body></html>`
	case "/myproblem":
		status = response.StatusInternalError
		html = `<html><body><h1>Internal Server Error</h1></body></html>`
	default:
		status = response.StatusOK
		html = `<html><body><h1>Success!</h1></body></html>`
	}

	headers := response.GetDefaultHeaders(len(html))
	headers.Set("Content-Type", "text/html")
	w.WriteStatusLine(status)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(html))
}
func main() {
	srv, err := server.Serve(port, myHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

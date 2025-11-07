package main

import (
	"github.com/RayanMalki/tcptohttp/internal/request"
	"github.com/RayanMalki/tcptohttp/internal/response"
	"github.com/RayanMalki/tcptohttp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func myHandler(w *response.Writer, req *request.Request) {
	var html string
	var status response.StatusCode

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.StatusBadRequest
		html = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
	case "/myproblem":
		status = response.StatusInternalError
		html = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`
	default:
		status = response.StatusOK
		html = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
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

package server

import (
	"fmt"
	"io"
	"net"

	"github.com/RayanMalki/tcptohttp/internal/request"
	"github.com/RayanMalki/tcptohttp/internal/response"
)

type Server struct {
	closed bool
}

type Handler func(w *response.Writer, req *request.Request)

func runConnection(s *Server, conn io.ReadWriteCloser, handler Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		w := response.NewWriter(conn)
		html := `<html>
  <head><title>400 Bad Request</title></head>
  <body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body>
</html>`
		w.WriteStatusLine(response.StatusBadRequest)
		headers := response.GetDefaultHeaders(len(html))
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody([]byte(html))
		return
	}

	w := response.NewWriter(conn)
	handler(w, req)
}

func runServer(s *Server, listener net.Listener, handler Handler) {
	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}
		go runConnection(s, conn, handler)
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{}
	go runServer(server, listener, handler)
	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}

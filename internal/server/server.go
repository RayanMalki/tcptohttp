package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/RayanMalki/tcptohttp/internal/request"
	"github.com/RayanMalki/tcptohttp/internal/response"
)

type Server struct {
	closed bool
}

type HandlerError struct {
	StatusCode     int
	HandlerMessage string
}

func runConnection(s *Server, conn io.ReadWriteCloser, handler Handler) {
	defer conn.Close()

	// Step 1: Parse the incoming request
	req, err := request.RequestFromReader(conn)
	if err != nil {
		writeErr := HandlerError{
			StatusCode:     400,
			HandlerMessage: "Bad Request\n",
		}
		WriteHandlerError(conn, writeErr)
		return
	}

	// Step 2: Create a new bytes.Buffer for the handler to write to
	var buffer bytes.Buffer

	// Step 3: Call the handler
	handlerErr := handler(&buffer, req)

	// Step 4: If handler returned an error, write it to the connection
	if handlerErr != nil {
		WriteHandlerError(conn, *handlerErr)
		return
	}

	// Step 5: Otherwise, send a normal 200 OK response
	body := buffer.Bytes()

	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		fmt.Println("error writing status line:", err)
		return
	}

	headers := response.GetDefaultHeaders(len(body))
	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Println("error writing headers:", err)
		return
	}

	if _, err := conn.Write(body); err != nil {
		fmt.Println("error writing body:", err)
		return
	}
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, handlerError HandlerError) error {
	body := []byte(handlerError.HandlerMessage)

	if err := response.WriteStatusLine(w, response.StatusCode(handlerError.StatusCode)); err != nil {

		return err
	}

	headersMap := response.GetDefaultHeaders(len(body))

	if err := response.WriteHeaders(w, headersMap); err != nil {

		return err
	}

	if _, err := w.Write(body); err != nil {

		return err
	}

	return nil

}

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/RayanMalki/tcptohttp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %s", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Bad request from %s: %s", conn.RemoteAddr(), err)
		return
	}

	fmt.Println("Request line:")
	fmt.Println("- Method:", req.RequestLine.Method)
	fmt.Println("- Target:", req.RequestLine.RequestTarget)
	fmt.Println("- Version:", req.RequestLine.HttpVersion)

	fmt.Println("Headers:")
	for k, v := range req.Headers {
		fmt.Printf("  %s: %s\n", k, v)
	}

	if len(req.Body) > 0 {
		fmt.Println("Body:")
		fmt.Println(string(req.Body))
	} else {
		fmt.Println("Body: (empty)")
	}

	fmt.Println("Connection to", conn.RemoteAddr(), "closed")
}

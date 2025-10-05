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
			log.Fatalf("error: %s", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("Bad request from %s: %s", conn.RemoteAddr(), err)
			continue
		}

		fmt.Println("Request line:")
		fmt.Println("-Method: " + req.RequestLine.Method)
		fmt.Println("-Target: " + req.RequestLine.RequestTarget)
		fmt.Println("-Version: " + req.RequestLine.HttpVersion)

		fmt.Println("Connection to", conn.RemoteAddr(), "closed")

	}

}


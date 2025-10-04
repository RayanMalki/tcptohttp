package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		str := ""

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)

			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			data = data[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				ch <- str
				str = ""
				data = data[i+1:]
			}

			str += string(data)
		}

		if len(str) != 0 {
			ch <- str
		}
	}()

	return ch
}

const port = ":42069"

func main() {

	listener, err := net.Listen("tcp4", port)

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

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println(line)
		}

		fmt.Println("Connection to", conn.RemoteAddr(), "closed")

	}

}

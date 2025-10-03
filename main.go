package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		log.Fatal("error", err)
	}

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}

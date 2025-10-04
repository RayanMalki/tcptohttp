package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	serverAddr := "localhost:42069"

	udpaddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP address: %v\n", err)
		os.Exit(1)

	}

	conn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing UDP: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Type you message and press Enter to send. Press Ctrl+C to exit.\n", serverAddr)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error while reading input: %v\n", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)

		}

		fmt.Printf("Message sent: %s", line)

	}

}

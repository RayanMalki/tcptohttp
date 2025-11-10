### HTTP/1.1 Server in Go

(RFC 9110 & RFC 9112 compliant)

## English Version

## Overview

This project implements a low-level HTTP/1.1 server written entirely in Go using raw TCP sockets, without relying on Go’s standard net/http package.
It faithfully follows the specifications of RFC 9110 (HTTP Semantics) and RFC 9112 (HTTP/1.1 Protocol), making it a standards-compliant and educational implementation.

The server supports:
	•	Full HTTP/1.1 request parsing (request line, headers, body, status lines)
	•	Chunked transfer encoding and trailers
	•	Proxy streaming via httpbin.org for real-time data forwarding
	•	Binary data handling (including MP4 streaming)
	•	Modular architecture ready for backend and frontend integration

  <img width="1536" height="1024" alt="architecture " src="https://github.com/user-attachments/assets/cf577d18-15ea-4e5f-89da-dce89911cb71" />

## Features:

HTTP Parsing: Implements request and response handling as per RFC 9110 and 9112

Chunked Transfer Encoding: Supports streaming data in hexadecimal-sized chunks

Trailers: Adds X-Content-SHA256 and X-Content-Length for data integrity verification

Proxy Mode: /httpbin/... routes forward live HTTP streams from https://httpbin.org

Binary Streaming: Streams binary files such as MP4 videos using proper content types

Extensibility: Designed for integration with Node.js backend and React frontend

Usage:

## Run the server
```bash
go run cmd/httpserver/main.go
```
## Test Proxy streaming (view raw chunked response)
```bash
echo -e "GET /httpbin/stream/100 HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" | nc localhost 42069
```
## Test binary data straming
In your browser paste the following:
```bash
http://localhost:42069/video/httptotcp.mp4
```

## Author
Rayan Malki -
Software Engineering Student

École de technologie supérieure (ÉTS), Montréal

Passionate about backend systems, networking, and low-level programming.










# HTTP/1.1 Server over Raw TCP (Go)

## Overview

This project implements a low-level HTTP/1.1 server written in Go using **raw TCP sockets**, without relying on Go’s standard `net/http` package.

It closely follows:

- **RFC 9110** – HTTP Semantics  
- **RFC 9112** – HTTP/1.1 Protocol  

The goal is to provide a **standards-compliant** and **educational** implementation of HTTP built directly on top of TCP.

---

## Features

- **Full HTTP/1.1 request parsing**
  - Request line, headers, body
  - Response/status line handling

- **Chunked transfer encoding**
  - Supports streaming data using hexadecimal-sized chunks
  - Proper handling of chunked bodies and terminating chunk

- **Trailers for data integrity**
  - Adds `X-Content-SHA256` and `X-Content-Length` as HTTP trailers
  - Enables verification of streamed content

- **Proxy streaming mode**
  - Routes under `/httpbin/...` forward live HTTP streams from `https://httpbin.org`
  - Acts as a basic HTTP proxy on top of the custom TCP server

- **Binary data streaming**
  - Streams binary files (e.g. MP4 videos)
  - Uses correct `Content-Type` and `Content-Length` headers

- **Extensible architecture**
  - Designed to integrate with a Node.js/TypeScript backend and a React frontend
  - Can be used as the “data plane” for monitoring/observability tools

---

## Architecture 

- **Transport layer:** Raw TCP socket listener in Go  
- **Protocol layer:** HTTP/1.1 parsing (request line, headers, body, chunked encoding)  
- **Application layer:**
  - Static/binary file serving
  - Proxy routes (`/httpbin/...`)
  - Trailers and integrity metadata

 ```mermaid
flowchart LR
    %% Client
    subgraph CLIENT[Client Side]
        C[Client (curl / browser)]
    end

    %% Go HTTP server
    subgraph GO[Go HTTP Server]
        M[Server - main.go]
        RP[Request Parser<br/>package request]
        RW[Response Writer<br/>package response]
        H[Handler Logic<br/>myHandler]
        M --> RP --> RW --> H
    end

    %% External service
    subgraph EXT[External Service]
        E[(httpbin.org)]
    end

    %% Flows
    C -- "TCP Request (HTTP/1.1)" --> M

    H -- "/httpbin/stream/x<br/>Proxy request" --> E
    E -- "Chunked Body + Trailers" --> RW
    RW -- "Chunked Body + Trailers" --> C

```


## Usage

### Run the server

```bash
go run cmd/httpserver/main.go
```

By default, the server listens on:
```bash
localhost:42069
```

### Test proxy streaming (raw chunked response)


```bash
echo -e "GET /httpbin/stream/100 HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" \
  | nc localhost 42069
```

### Test binary data streaming (MP4)

Open your browser and navigate to:
```bash
http://localhost:42069/video/httptotcp.mp4
```

### Author

Rayan Malki
Software Engineering Student – École de technologie supérieure (ÉTS), Montréal

Passionate about:
	•	Backend systems and distributed architectures
	•	Networking and protocol design
	•	Low-level programming and infrastructure tooling





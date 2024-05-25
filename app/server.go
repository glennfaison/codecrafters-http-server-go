package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	ReadFromConnection(conn)
}

func ReadFromConnection(conn net.Conn) {
	const Http_Status_200 = "HTTP/1.1 200 OK\r\n\r\n"
	const Http_Status_404 = "HTTP/1.1 404 Not Found\r\n\r\n"

	buff := make([]byte, 1024)
	if _, err := conn.Read(buff); err != nil {
		fmt.Println("Failed to read from connection")
		os.Exit(1)
	}

	if strings.HasPrefix(string(buff), "GET / HTTP/1.1") {
		conn.Write([]byte(Http_Status_200))
	} else {
		conn.Write([]byte(Http_Status_404))
	}
}

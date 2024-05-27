package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening...")

	// connectionChannel := make(chan net.Conn)
	for {
		fmt.Println("before accepting connection")
		connection, err := listener.Accept()
		fmt.Println("Connection accepted")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	const Http_Status_200 = "HTTP/1.1 200 OK\r\n\r\n"
	const Http_Status_404 = "HTTP/1.1 404 Not Found\r\n\r\n"
	const Http_Status_500 = "HTTP/1.1 500 Internal Server Error\r\n\r\n"

	buff := make([]byte, 1024)
	charactersRead, err := connection.Read(buff)
	if err != nil {
		fmt.Println("Failed to read from connection")
		connection.Write([]byte(Http_Status_500))
	}
	stringReceived := string(buff[:charactersRead])
	println("received string", stringReceived)

	pattern, err := regexp.Compile(`GET /echo/(\w+)`)
	if err != nil {
		fmt.Println("Failed to parse request")
	}

	strs := pattern.FindStringSubmatch(stringReceived)
	if len(strs) > 0 {
		match := strs[1]
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(match)) + "\r\n\r\n" + match
		connection.Write([]byte(response))
		return
	}

	if strings.HasPrefix(string(buff), "GET / HTTP/1.1") {
		connection.Write([]byte(Http_Status_200))
	} else {
		connection.Write([]byte(Http_Status_404))
	}
}

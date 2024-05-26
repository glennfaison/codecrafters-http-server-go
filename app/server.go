package main

import (
	"fmt"
	"regexp"
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
	fmt.Println("Listening...")

	connectionChannel := make(chan net.Conn)
	for {
		fmt.Println("before accepting connection")
		conn, err := listener.Accept()
		fmt.Println("Connection accepted")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Sending connection to channel...")
		connectionChannel <- conn
		fmt.Println("Connection sent to channel...")
		<-connectionChannel
		go ReadFromConnection(connectionChannel)
	}
}

func ReadFromConnection(connectionChannel chan net.Conn) {
	fmt.Println("reading from connection...")
	conn := <-connectionChannel
	const Http_Status_200 = "HTTP/1.1 200 OK\r\n\r\n"
	const Http_Status_404 = "HTTP/1.1 404 Not Found\r\n\r\n"

	buff := make([]byte, 1024)
	if _, err := conn.Read(buff); err != nil {
		fmt.Println("Failed to read from connection")
		os.Exit(1)
	}

	pattern, err := regexp.Compile(`GET /echo/(\w+)`)
	if err != nil {
		fmt.Println("Failed to parse request")
	}
	strs := pattern.FindAllString(string(buff), -1)
	for _, str := range strs {
		println(str)
	}

	if strings.HasPrefix(string(buff), "GET / HTTP/1.1") {
		conn.Write([]byte(Http_Status_200))
	} else {
		conn.Write([]byte(Http_Status_404))
	}
}

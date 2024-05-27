package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	myhttp "github.com/codecrafters-io/http-server-starter-go/app/pkg/my-http"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening...")

	for {
		fmt.Println("before accepting connection")
		connection, err := listener.Accept()
		fmt.Println("Connection accepted")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	request, err := myhttp.ParseRequest(connection)
	if err != nil {
		fmt.Println("Failed to read from connection")
		connection.Write([]byte(myhttp.NewResponse().SetStatus(500).ToString()))
		return
	}

	switch true {
	case request.GetPath() == "/":
		connection.Write([]byte(myhttp.NewResponse().SetStatus(200).ToString()))
	case strings.HasPrefix(request.GetPath(), "/user-agent"):
		response := myhttp.NewResponse().SetStatus(200).SetBody(request.GetHeader("user-agent")).ToString()
		println(response)
		connection.Write([]byte(response))
	case strings.HasPrefix(request.GetPath(), "/echo/"):
		pattern, err := regexp.Compile(`/echo/(\w*)`)
		if err != nil {
			fmt.Println("Failed to parse request path")
		}

		strs := pattern.FindStringSubmatch(request.GetPath())
		if len(strs) > 0 {
			match := strs[1]
			response := myhttp.NewResponse().SetStatus(200).SetBody(match).ToString()
			connection.Write([]byte(response))
		}
	default:
		connection.Write([]byte(myhttp.NewResponse().SetStatus(404).ToString()))
	}
}

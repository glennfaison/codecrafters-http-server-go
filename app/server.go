package main

import (
	"encoding/json"
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

	request, err := parseRequest(connection)
	if err != nil {
		fmt.Println("Failed to read from connection")
		connection.Write([]byte(Http_Status_500))
		return
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		fmt.Println("failed to marshal request data to JSON")
	}
	println(requestJson)

	switch true {
	case request.Path == "/":
		connection.Write([]byte(Http_Status_200))
	case strings.HasPrefix(request.Path, "/user-agent"):
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(request.Headers["user-agent"])) + "\r\n\r\n" + request.Headers["User-Agent"]
		connection.Write([]byte(response))
	case strings.HasPrefix(request.Path, "/echo/"):
		pattern, err := regexp.Compile(`/echo/(\w*)`)
		if err != nil {
			fmt.Println("Failed to parse request")
		}

		strs := pattern.FindStringSubmatch(request.Path)
		if len(strs) > 0 {
			match := strs[1]
			response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(match)) + "\r\n\r\n" + match
			connection.Write([]byte(response))
		}
	default:
		connection.Write([]byte(Http_Status_404))
	}
}

type MyHttpRequest struct {
	Method   string            `json:"method"`
	Path     string            `json:"path"`
	Protocol string            `json:"protocol"`
	Headers  map[string]string `json:"Headers"`
}

func parseRequest(connection net.Conn) (MyHttpRequest, error) {
	buff := make([]byte, 1024)
	charactersRead, err := connection.Read(buff)
	if err != nil {
		fmt.Println("Failed to read from connection")
		return MyHttpRequest{}, err
	}
	requestString := string(buff[:charactersRead])
	println("parseRequest: read", requestString)

	httpRequestLineRegexp, err := regexp.Compile(`([A-Z]+) (\S+) (\S+)`)
	if err != nil {
		fmt.Println("Failed to parse request line")
		return MyHttpRequest{}, err
	}
	requestLineMatches := httpRequestLineRegexp.FindStringSubmatch(requestString)

	requestHeaders, err := parseRequestHeaders(requestString)
	if err != nil {
		return MyHttpRequest{}, err
	}

	var response MyHttpRequest = MyHttpRequest{
		Method:   requestLineMatches[1],
		Path:     requestLineMatches[2],
		Protocol: requestLineMatches[3],
		Headers:  requestHeaders,
	}

	return response, nil
}

func parseRequestHeaders(requestString string) (map[string]string, error) {
	httpRequestHeadersRegexp, err := regexp.Compile(`(\s+((\w|-)+): (\S+))`)
	if err != nil {
		fmt.Println("Failed to parse request headers")
		return nil, err
	}
	requestHeadersMatches := httpRequestHeadersRegexp.FindAllStringSubmatch(requestString, -1)
	requestHeaders := map[string]string{}
	for i := 0; i < len(requestHeadersMatches); i++ {
		requestHeaders[strings.ToLower(requestHeadersMatches[i][2])] = requestHeadersMatches[i][4]
	}
	return requestHeaders, nil
}

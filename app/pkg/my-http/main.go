package myhttp

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func ParseRequest(connection net.Conn) (Request, error) {
	buff := make([]byte, 1024)
	charactersRead, err := connection.Read(buff)
	if err != nil {
		fmt.Println("Failed to read from connection")
		return Request{}, err
	}
	requestString := string(buff[:charactersRead])

	// Parse request line.
	httpRequestLineRegexp, err := regexp.Compile(`([A-Z]+) (\S+) (\S+)`)
	if err != nil {
		fmt.Println("Failed to parse request line")
		return Request{}, err
	}
	requestLineMatches := httpRequestLineRegexp.FindStringSubmatch(requestString)

	// Parse request headers.
	requestHeaders, err := ParseRequestHeaders(requestString)
	if err != nil {
		return Request{}, err
	}

	// Parse request body.
	requestBodyRegexp, err := regexp.Compile(`\r\n\r\n(.*)`)
	if err != nil {
		fmt.Println("Failed to parse request body")
		return Request{}, err
	}
	requestBodyMatches := requestBodyRegexp.FindStringSubmatch(requestString)
	requestBodyString := ""
	if len(requestBodyMatches) >= 2 {
		requestBodyString = requestBodyMatches[1]
	}

	var response Request = Request{
		method:   requestLineMatches[1],
		path:     requestLineMatches[2],
		protocol: requestLineMatches[3],
		headers:  requestHeaders,
		body:     requestBodyString,
	}

	return response, nil
}

func ParseRequestHeaders(requestString string) (map[string]string, error) {
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

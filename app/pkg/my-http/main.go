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
	println("parseRequest: read", requestString)

	httpRequestLineRegexp, err := regexp.Compile(`([A-Z]+) (\S+) (\S+)`)
	if err != nil {
		fmt.Println("Failed to parse request line")
		return Request{}, err
	}
	requestLineMatches := httpRequestLineRegexp.FindStringSubmatch(requestString)

	requestHeaders, err := ParseRequestHeaders(requestString)
	if err != nil {
		return Request{}, err
	}

	var response Request = Request{
		method:   requestLineMatches[1],
		path:     requestLineMatches[2],
		protocol: requestLineMatches[3],
		headers:  requestHeaders,
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

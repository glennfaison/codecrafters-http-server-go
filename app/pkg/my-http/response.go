package myhttp

import (
	"strconv"
	"strings"
)

type Response struct {
	protocol string
	status   int
	headers  map[string]string
	body     string
}

func NewResponse() *Response {
	return &Response{
		protocol: "HTTP/1.1",
		headers: map[string]string{
			"content-type": "text/plain",
		},
	}
}

func (r *Response) SetProtocol(protocol string) *Response {
	r.protocol = protocol
	return r
}

func (r *Response) SetStatus(status int) *Response {
	r.status = status
	return r
}

func (r *Response) SetHeaders(headers map[string]string) *Response {
	for name, value := range headers {
		r.headers[strings.ToLower(name)] = value
	}
	return r
}

func (r *Response) AddHeader(name string, value string) *Response {
	r.headers[strings.ToLower(name)] = value
	return r
}

func (r *Response) RemoveHeader(name string) *Response {
	delete(r.headers, strings.ToLower(name))
	return r
}

func (r *Response) SetBody(body string) *Response {
	r.body = body
	r.AddHeader("content-length", strconv.Itoa(len(body)))
	return r
}

func (r *Response) ToString() string {
	Http_Status_Codes := map[int]string{}
	Http_Status_Codes[200] = "200 OK\r\n\r\n"
	Http_Status_Codes[404] = "404 Not Found\r\n\r\n"
	Http_Status_Codes[500] = "500 Internal Server Error\r\n\r\n"
	statusString, found := Http_Status_Codes[r.status]
	if !found {
		statusString = Http_Status_Codes[500]
	}

	// Add response status line.
	responseString := r.protocol + " " + statusString
	// Add response headers, each on a line.
	for name, value := range r.headers {
		responseString += name + ": " + value + "\r\n"
	}
	responseString += "\r\n"
	// Add response body
	responseString += r.body
	return responseString
}

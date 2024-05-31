package myhttp

import "strings"

type Request struct {
	method   string
	path     string
	protocol string
	headers  map[string]string
	params   map[string]string
	query    map[string]string
	body     string
}

func (r Request) GetMethod() string {
	return r.method
}

func (r Request) GetPath() string {
	return r.path
}

func (r Request) GetProtocol() string {
	return r.protocol
}

func (r Request) GetHeaders() map[string]string {
	return r.headers
}

func (r Request) GetHeader(header string) string {
	return r.headers[strings.ToLower(header)]
}

func (r Request) HasHeader(header string) bool {
	_, exists := r.headers[strings.ToLower(header)]
	return exists
}

func (r Request) GetBody() string {
	return r.body
}

func (r Request) GetParams() map[string]string {
	return r.params
}

func (r Request) GetParam(name string) string {
	return r.params[name]
}

func (r Request) AddParam(name string, value string) {
	r.params[name] = value
}

func (r Request) GetQueryParams() map[string]string {
	return r.query
}

func (r Request) GetQueryParam(name string) string {
	return r.query[name]
}

func (r Request) AddQueryParam(name string, value string) {
	r.query[name] = value
}

func (r *Request) ToString() string {
	requestString := r.method + " " + r.path + " " + r.protocol + "\r\n"
	for name, value := range r.headers {
		requestString += name + ": " + value + "\r\n"
	}
	requestString += "\r\n"
	requestString += r.body
	return requestString
}

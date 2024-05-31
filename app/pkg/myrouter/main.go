package myrouter

import (
	"fmt"
	"net"
	"strings"

	myhttp "github.com/codecrafters-io/http-server-starter-go/app/pkg/my-http"
)

type Node struct {
	name     string
	handler  func(net.Conn, myhttp.Request)
	children map[string]*Node
}

type Router struct {
	listener           net.Listener
	head               *Node
	fallThroughHandler func(net.Conn, myhttp.Request)
}

func NewNode() *Node {
	return &Node{
		name:     "",
		children: map[string]*Node{},
	}
}

func NewRouter(listener net.Listener) *Router {
	return &Router{
		listener: listener,
		head:     NewNode(),
	}
}

func (t *Router) GetCallbackKeys(method string, urlPattern string) []string {
	queryIndex := strings.Index(urlPattern, "?")
	hasQuery := queryIndex >= 0
	newUrlPattern := urlPattern
	if hasQuery {
		newUrlPattern = urlPattern[:queryIndex]
	}
	if !strings.HasPrefix(newUrlPattern, "/") {
		newUrlPattern = "/" + newUrlPattern
	}
	patternSections := strings.Split(strings.ToLower(newUrlPattern), "/")
	patternSections = patternSections[1:]
	return append([]string{strings.ToLower(method)}, patternSections...)
}

func (t *Router) RegisterRouteHandler(method string, urlPattern string, handler func(net.Conn, myhttp.Request)) {
	callbackKeys := t.GetCallbackKeys(method, urlPattern)

	ptr := t.head
	println(method, urlPattern)
	fmt.Printf("%v %d\n", callbackKeys, len(callbackKeys))
	for _, key := range callbackKeys {
		name := key
		if strings.HasPrefix(key, ":") {
			key = ":var"
		}
		println("LOG:", key, ptr.name == "", len(ptr.children))
		if _, exists := ptr.children[key]; !exists {
			ptr.children[key] = NewNode()
			ptr.children[key].name = name
		}
		ptr = ptr.children[key]
	}
	ptr.handler = handler
}

func (t *Router) RegisterFallthroughHandler(handler func(net.Conn, myhttp.Request)) {
	t.fallThroughHandler = handler
}

func (t Router) GetHandler(method string, urlPattern string) func(net.Conn, myhttp.Request) {
	callbackKeys := t.GetCallbackKeys(method, urlPattern)

	ptr := t.head
	for _, key := range callbackKeys {
		if _, exists := ptr.children[key]; !exists {
			key = ":var"
		}
		ptr = ptr.children[key]
		if ptr == nil {
			break
		}
	}
	if ptr == nil || ptr.handler == nil {
		return t.fallThroughHandler
	}
	return ptr.handler
}

func (r Router) handleConnection(connection net.Conn) {
	defer connection.Close()

	request, err := myhttp.ParseRequest(connection)
	if err != nil {
		fmt.Println("Failed to read from connection")
		response := myhttp.NewResponse().SetStatus(500)
		connection.Write(response.ToBytes())
		return
	}

	handlerFn := r.GetHandler(request.GetMethod(), request.GetPath())
	handlerFn(connection, request)
}

func (r Router) Start() {
	for {
		connection, err := r.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go r.handleConnection(connection)
	}
}

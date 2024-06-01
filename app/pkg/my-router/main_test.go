package myrouter

import (
	"net"
	"reflect"
	"testing"

	myhttp "github.com/codecrafters-io/http-server-starter-go/app/pkg/my-http"
)

func TestRouter_Get_AND_Register_RouteHandler(t *testing.T) {
	handlerFn := func(c net.Conn, r myhttp.Request) { print("handler") }
	fallThroughHandlerFn := func(c net.Conn, r myhttp.Request) { print("fallThroughHandlerFn") }

	type args struct {
		method          string
		urlPattern      string
		url             string
		handler         func(net.Conn, myhttp.Request)
		expectedHandler func(net.Conn, myhttp.Request)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GET /user-agent HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			args: args{
				method:          "GET",
				urlPattern:      "/user-agent",
				url:             "user-agent",
				handler:         handlerFn,
				expectedHandler: handlerFn,
			},
		},
		{
			name: "GET /user-agent HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			args: args{
				method:          "GET",
				urlPattern:      "/user-agent",
				url:             "/user-agent",
				handler:         handlerFn,
				expectedHandler: handlerFn,
			},
		},
		{
			name: "GET /user-agent HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			args: args{
				method:          "GET",
				urlPattern:      "/user-agent",
				url:             "/user-agent/",
				handler:         handlerFn,
				expectedHandler: handlerFn,
			},
		},
		{
			name: "GET /user-agent HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			args: args{
				method:          "GET",
				urlPattern:      "/user-agent",
				url:             "/user-agents/",
				handler:         handlerFn,
				expectedHandler: fallThroughHandlerFn,
			},
		},
		{
			name: "GET /user-agent HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			args: args{
				method:          "GET",
				urlPattern:      "/user-agent",
				url:             "/user-agent/test",
				handler:         handlerFn,
				expectedHandler: fallThroughHandlerFn,
			},
		},
	}
	router := Router{
		listener:           &net.TCPListener{},
		head:               NewNode(),
		fallThroughHandler: fallThroughHandlerFn,
	}
	for _, tt := range tests {
		router.RegisterRouteHandler(tt.args.method, tt.args.url, tt.args.handler)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := router.GetRouteHandler(tt.args.method, tt.args.urlPattern)
			if reflect.DeepEqual(got, tt.args.expectedHandler) {
				t.Errorf("Router.GetRouteHandler() result doesn't match expected handler")
			}
		})
	}
}

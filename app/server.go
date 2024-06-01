package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"

	myhttp "github.com/codecrafters-io/http-server-starter-go/app/pkg/my-http"
	myexpress "github.com/codecrafters-io/http-server-starter-go/app/pkg/my-router"
)

var directory string

func main() {
	flag.StringVar(&directory, "directory", "", "Absolute path to the server file.")
	flag.Parse()
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening...")

	router := myexpress.NewRouter(listener)
	router.RegisterRouteHandler("GET", "/", func(connection net.Conn, request myhttp.Request) {
		response := myhttp.NewResponse().SetStatus(200).ToBytes()
		connection.Write(response)
	})
	router.RegisterRouteHandler("GET", "/user-agent", func(connection net.Conn, request myhttp.Request) {
		response := myhttp.NewResponse().SetStatus(200).SetBody(request.GetHeader("user-agent")).ToBytes()
		connection.Write(response)
	})
	router.RegisterRouteHandler("GET", "/echo/:value", func(connection net.Conn, request myhttp.Request) {
		pattern, err := regexp.Compile(`/echo/(\w*)`)
		if err != nil {
			fmt.Println("Failed to parse request path")
		}

		strs := pattern.FindStringSubmatch(request.GetPath())
		if len(strs) > 0 {
			match := strs[1]
			response := myhttp.NewResponse().SetStatus(200).SetBody(match)
			x := strings.Split(request.GetHeader("accept-encoding"), ",")
			for idx := range x {
				x[idx] = strings.TrimSpace(x[idx])
			}
			if slices.Index(x, "gzip") >= 0 {
				response.AddHeader("Content-Encoding", "gzip")
				var b bytes.Buffer
				gz := gzip.NewWriter(&b)
				gz.Close() // TODO: Find out why it MUST be closed immediately in order for the result in `b` to make sense
				gz.Write([]byte(match))
				response.SetBody(b.String())
			}
			connection.Write([]byte(response.ToString()))
		}
	})
	router.RegisterFallthroughHandler(func(connection net.Conn, request myhttp.Request) {
		re, err := regexp.Compile(`/files/(\S+)`)
		if err != nil {
			fmt.Println("Failed to parse request path")
			return
		}
		matches := re.FindStringSubmatch(request.GetPath())
		if len(matches) == 2 {
			filePath := path.Join(directory, matches[1])
			if request.GetMethod() == http.MethodGet {
				if err := getFile(connection, filePath); err != nil {
					return
				}
			}
			if request.GetMethod() == http.MethodPost {
				if err := postFile(request, connection, filePath); err != nil {
					return
				}
			}
		}

		response := myhttp.NewResponse().SetStatus(404).ToBytes()
		connection.Write(response)
	})

	router.Start()
}

func getFile(connection net.Conn, filePath string) error {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		responseStr := myhttp.NewResponse().SetStatus(404).ToString()
		connection.Write([]byte(responseStr))
		return err
	}
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read file: ", filePath, err.Error())
		return err
	}
	responseStr := myhttp.NewResponse().SetStatus(200).AddHeader("content-type", "application/octet-stream").SetBody(string(fileData)).ToString()
	connection.Write([]byte(responseStr))
	return nil
}

func postFile(request myhttp.Request, connection net.Conn, filePath string) error {
	x := request.GetBody()
	if err := os.WriteFile(filePath, []byte(x), 0644); err != nil {
		return err
	}
	responseStr := myhttp.NewResponse().SetStatus(201).ToString()
	connection.Write([]byte(responseStr))
	return nil
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	directoryPtr := flag.String("directory", "empty", "Set the current directory of the server")

	flag.Parse()

	server := NewServer(*directoryPtr)

	server.Get("/", func(request HTTPRequest) HTTPResponse {
		resp := NewHTTPResponse(OK)

		data := "<h1>Hello World</h1>"

		resp.SetHeader("Content-Type", "text/html")
		resp.SetHeader("Content-Length", fmt.Sprint(len(data)))
		resp.SetBody(data)

		return resp
	})

	server.Get("/test", func(request HTTPRequest) HTTPResponse {
		resp := NewHTTPResponse(OK)

		data := `<h1 style="color: blue;">Hello From Test</h1><button style="color: red; background-color: black;">Click me</button>`

		resp.SetHeader("Content-Type", "text/html")
		resp.SetHeader("Content-Length", fmt.Sprint(len(data)))
		resp.SetBody(data)

		return resp
	})

	server.Get("/user-agent", func(request HTTPRequest) HTTPResponse {
		response := NewHTTPResponse(OK)

		response.SetHeader("Content-Type", "text/plain")
		response.SetHeader("Content-Length", fmt.Sprint(len(request.UserAgent)))
		response.SetBody(request.UserAgent)

		return response
	})

	server.Get("/echo/*", func(request HTTPRequest) HTTPResponse {
		echo, _ := strings.CutPrefix(request.Path, "/echo/")

		response := NewHTTPResponse(OK)

		response.SetHeader("Content-Type", "text/plain")
		response.SetHeader("Content-Length", fmt.Sprint(len(echo)))
		response.SetBody(echo)

		return response
	})

	server.Get("/files/*", func(request HTTPRequest) HTTPResponse {
		fileName, _ := strings.CutPrefix(request.Path, "/files/")
		fullPath := server.Directory + fileName

		_, err := os.Stat(fullPath)

		if os.IsNotExist(err) {
			return NewHTTPResponse(NOTFOUND)
		}

		response := NewHTTPResponse(OK)

		file, _ := os.Open(fullPath)
		defer file.Close()

		fileContents, _ := io.ReadAll(file)
		fileContentsStr := string(fileContents)

		response.SetHeader("Content-Type", "application/octet-stream")
		response.SetHeader("Content-Length", fmt.Sprint(len(fileContentsStr)))
		response.SetBody(fileContentsStr)

		return response
	})

	server.Post("/files/*", func(request HTTPRequest) HTTPResponse {
		fileName, _ := strings.CutPrefix(request.Path, "/files/")
		fullPath := server.Directory + fileName

		response := NewHTTPResponse(CREATED)

		file, _ := os.Create(fullPath)
		defer file.Close()

		data := []byte(request.Body)
		file.Write(data)

		return response
	})

	server.Get("/basketball/party/*", func(request HTTPRequest) HTTPResponse {
		response := NewHTTPResponse(OK)

		response.SetBody("Let's play basketball and party")

		return response
	})

	err := server.Start()
	if err != nil {
		slog.Error("There was an issue running the server: ", err)
	}
}

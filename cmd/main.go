package main

import (
	"fmt"
	"github.com/godofprodev/http-server/internal/config"
	"github.com/godofprodev/http-server/internal/server"
	"io"
	"log/slog"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	cfg := config.NewConfig()

	srv := server.NewServer(cfg)

	srv.Get("/", func(request server.HTTPRequest) server.HTTPResponse {
		resp := server.NewHTTPResponse(server.OK)

		data := "<h1>Hello World</h1>"

		resp.SetHeader("Content-Type", "text/html")
		resp.SetHeader("Content-Length", fmt.Sprint(len(data)))
		resp.SetBody(data)

		return resp
	})

	srv.Get("/test", func(request server.HTTPRequest) server.HTTPResponse {
		resp := server.NewHTTPResponse(server.OK)

		data := `<h1 style="color: blue;">Hello From Test</h1><button style="color: red; background-color: black;">Click me</button>`

		resp.SetHeader("Content-Type", "text/html")
		resp.SetHeader("Content-Length", fmt.Sprint(len(data)))
		resp.SetBody(data)

		return resp
	})

	srv.Get("/user-agent", func(request server.HTTPRequest) server.HTTPResponse {
		response := server.NewHTTPResponse(server.OK)

		response.SetHeader("Content-Type", "text/plain")
		response.SetHeader("Content-Length", fmt.Sprint(len(request.UserAgent)))
		response.SetBody(request.UserAgent)

		return response
	})

	srv.Get("/echo/*", func(request server.HTTPRequest) server.HTTPResponse {
		echo, _ := strings.CutPrefix(request.Path, "/echo/")

		response := server.NewHTTPResponse(server.OK)

		response.SetHeader("Content-Type", "text/plain")
		response.SetHeader("Content-Length", fmt.Sprint(len(echo)))
		response.SetBody(echo)

		return response
	})

	srv.Get("/files/*", func(request server.HTTPRequest) server.HTTPResponse {
		fileName, _ := strings.CutPrefix(request.Path, "/files/")
		fullPath := srv.Config.Directory + fileName

		_, err := os.Stat(fullPath)

		if os.IsNotExist(err) {
			return server.NewHTTPResponse(server.NOTFOUND)
		}

		response := server.NewHTTPResponse(server.OK)

		file, _ := os.Open(fullPath)
		defer file.Close()

		fileContents, _ := io.ReadAll(file)
		fileContentsStr := string(fileContents)

		response.SetHeader("Content-Type", "application/octet-stream")
		response.SetHeader("Content-Length", fmt.Sprint(len(fileContentsStr)))
		response.SetBody(fileContentsStr)

		return response
	})

	srv.Post("/files/*", func(request server.HTTPRequest) server.HTTPResponse {
		fileName, _ := strings.CutPrefix(request.Path, "/files/")
		fullPath := srv.Config.Directory + fileName

		response := server.NewHTTPResponse(server.CREATED)

		file, _ := os.Create(fullPath)
		defer file.Close()

		data := []byte(request.Body)
		file.Write(data)

		return response
	})

	srv.Get("/basketball/party/*", func(request server.HTTPRequest) server.HTTPResponse {
		response := server.NewHTTPResponse(server.OK)

		response.SetBody("Let's play basketball and party")

		return response
	})

	err := srv.Start()
	if err != nil {
		slog.Error("There was an issue running the srv: ", err)
	}
}

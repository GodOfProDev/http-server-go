package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
)

type Server struct {
	Directory   string
	Listener    net.Listener
	GetPathMap  map[string]CallbackFunc
	PostPathMap map[string]CallbackFunc
}

func NewServer(directory string) Server {
	return Server{
		Directory:   directory,
		GetPathMap:  make(map[string]CallbackFunc),
		PostPathMap: make(map[string]CallbackFunc),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		return err
	}
	s.Listener = listener

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(conn)
	}
}

type CallbackFunc func(request HTTPRequest) HTTPResponse

func (s *Server) Get(path string, callback CallbackFunc) {
	s.GetPathMap[path] = callback
}

func (s *Server) Post(path string, callback CallbackFunc) {
	s.PostPathMap[path] = callback
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		slog.Error("Failed to read the connection: ", err)
		return
	}

	data := string(buf[:n])

	req := NewHTTPRequest(data)

	if req.Method == "GET" {
		var path string

		// TODO: MAKE WILDCARDS AND DON'T HARDCODE LIKE THIS
		if strings.Contains(req.Path, "/echo") {
			path = "/echo/*"
		} else if strings.Contains(req.Path, "/files") {
			path = "/files/*"
		} else {
			path = req.Path
		}

		callback, ok := s.GetPathMap[path]
		var response HTTPResponse
		if !ok {
			response = NewHTTPResponse(NOTFOUND)
		} else {
			response = callback(req)
		}

		writeToConnection(response.String(), conn)
	} else if req.Method == "POST" {
		var path string

		if strings.Contains(req.Path, "/files") {
			path = "/files/*"
		} else {
			path = req.Path
		}

		callback, ok := s.PostPathMap[path]
		var response HTTPResponse
		if !ok {
			response = NewHTTPResponse(NOTFOUND)
		} else {
			response = callback(req)
		}
		writeToConnection(response.String(), conn)
	}

	return

	header := "HTTP/1.1 "
	isEcho := false
	isUserAgent := false
	isFileGet := false
	isFilePost := false
	var fullFilePath string

	if req.Path == "/" {
		header += "200 OK\r\n\r\n"
	} else if strings.Contains(req.Path, "/echo") {
		header += "200 OK\r\n"
		isEcho = true

	} else if strings.Contains(req.Path, "/user-agent") {
		header += "200 OK\r\n"
		isUserAgent = true
	} else if strings.Contains(req.Path, "/files") && req.Method == "GET" {
		fileName, _ := strings.CutPrefix(req.Path, "/files/")
		fullFilePath = s.Directory + fileName

		_, err := os.Stat(fullFilePath)

		if os.IsNotExist(err) {
			header += "404 Not Found\r\n\r\n"
		} else {
			header += "200 OK\r\n"
			isFileGet = true
		}
	} else if strings.Contains(req.Path, "/files") && req.Method == "POST" {
		fileName, _ := strings.CutPrefix(req.Path, "/files/")
		fullFilePath = s.Directory + fileName

		header += "201 Created\r\n"
		isFilePost = true
	} else {
		header += "404 Not Found\r\n\r\n"
	}

	if strings.Contains(header, "\r\n\r\n") {
		writeToConnection(header, conn)
		return
	}

	if isEcho {
		strippedPath := stripPath(data)

		header += "Content-Type: text/plain\r\n"
		header += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(strippedPath))
		header += strippedPath + "\r\n\r\n"
	} else if isUserAgent {
		header += "Content-Type: text/plain\r\n"
		header += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(req.UserAgent))
		header += req.UserAgent
	} else if isFileGet {
		file, _ := os.Open(fullFilePath)
		defer file.Close()

		fileContents, _ := io.ReadAll(file)
		fileContentsStr := string(fileContents)

		header += "Content-Type: application/octet-stream\r\n"
		header += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(fileContentsStr))
		header += fileContentsStr
	} else if isFilePost {
		file, _ := os.Create(fullFilePath)
		defer file.Close()

		data := []byte(req.Body)
		file.Write(data)

		header += "\r\n"
	}

	writeToConnection(header, conn)
}

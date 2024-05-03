package server

import (
	"github.com/godofprodev/http-server/internal/config"
	"github.com/godofprodev/http-server/internal/util"
	"log/slog"
	"net"
	"regexp"
	"strings"
)

type Server struct {
	Config               *config.Config
	Listener             net.Listener
	GetPathMap           map[string]CallbackFunc
	PostPathMap          map[string]CallbackFunc
	GetPathWildcardsMap  map[*regexp.Regexp]CallbackFunc
	PostPathWildcardsMap map[*regexp.Regexp]CallbackFunc
}

func NewServer(cfg *config.Config) Server {
	return Server{
		Config:               cfg,
		GetPathMap:           make(map[string]CallbackFunc),
		PostPathMap:          make(map[string]CallbackFunc),
		GetPathWildcardsMap:  make(map[*regexp.Regexp]CallbackFunc),
		PostPathWildcardsMap: make(map[*regexp.Regexp]CallbackFunc),
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
	if strings.Contains(path, "/*") {
		str := strings.ReplaceAll(path, "/*", "(.*)")
		reg, _ := regexp.Compile(str)
		s.GetPathWildcardsMap[reg] = callback
	}

	s.GetPathMap[path] = callback
}

func (s *Server) Post(path string, callback CallbackFunc) {
	if strings.Contains(path, "/*") {
		str := strings.ReplaceAll(path, "/*", "(.*)")
		reg, _ := regexp.Compile(str)
		s.PostPathWildcardsMap[reg] = callback
	}

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

	callbackGet, okGet := s.GetPathMap[req.Path]
	callbackPost, okPost := s.PostPathMap[req.Path]

	response := NewHTTPResponse(NOTFOUND)

	if !okGet && okPost {
		response = callbackPost(req)
	} else if okGet && !okPost {
		response = callbackGet(req)
	}

	for r, callbackFunc := range s.GetPathWildcardsMap {
		matches := r.FindStringSubmatch(req.Path)

		if len(matches) > 1 {
			response = callbackFunc(req)
			break
		}
	}

	for r, callbackFunc := range s.PostPathWildcardsMap {
		matches := r.FindStringSubmatch(req.Path)

		if len(matches) > 1 {
			response = callbackFunc(req)
			break
		}
	}

	if response.StatusCode == NOTFOUND {
		response.SetBody("404 page not found")
	}

	util.WriteToConnection(response.String(), conn)

	return
}

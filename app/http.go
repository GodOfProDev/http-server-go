package main

import (
	"fmt"
	"strings"
)

const CRLF = "\r\n"

type HTTPResponse struct {
	StatusCode StatusCode
	Headers    map[string]string
	Body       string
}

type StatusCode int

func (s StatusCode) String() string {
	switch s {
	case OK:
		return "Ok"
	case CREATED:
		return "Created"
	case NOTFOUND:
		return "Not Found"
	}

	return ""
}

func (s StatusCode) Int() int {
	return int(s)
}

const (
	OK       StatusCode = 200
	CREATED  StatusCode = 201
	NOTFOUND StatusCode = 404
)

func NewHTTPResponse(statusCode StatusCode) HTTPResponse {
	return HTTPResponse{
		StatusCode: statusCode,
		Headers:    make(map[string]string),
	}
}

func (r *HTTPResponse) SetHeader(headerName string, headerContent string) {
	r.Headers[headerName] = headerContent
}

func (r *HTTPResponse) SetBody(body string) {
	r.Body = body
}

func (r *HTTPResponse) String() string {
	header := "HTTP/1.1 " + fmt.Sprint(r.StatusCode.Int()) + " " + r.StatusCode.String() + CRLF

	if len(r.Headers) == 0 {
		header += CRLF + CRLF

		if len(r.Body) != 0 {
			header += r.Body
		}

		return header
	}

	for k, v := range r.Headers {
		header += k + ": " + v + CRLF
	}

	header += CRLF

	if len(r.Body) != 0 {
		header += r.Body
	}

	return header
}

type HTTPRequest struct {
	Method    string
	Path      string
	Version   string
	UserAgent string
	Body      string
}

func NewHTTPRequest(header string) HTTPRequest {
	lines := strings.Split(header, "\r\n")
	headerMap := make(map[string]string)
	var body string
	var emptyLineIndex = -1

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if line == "" {
			emptyLineIndex = i
			break
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if i < emptyLineIndex || emptyLineIndex == -1 {
			if strings.Contains(line, "GET") || strings.Contains(line, "POST") {
				splitHeader := strings.Split(line, " ")
				headerMap["Method"] = splitHeader[0]
				headerMap["Path"] = splitHeader[1]
				headerMap["Version"] = splitHeader[2]
			}

			if strings.Contains(line, "User-Agent") {
				content, _ := strings.CutPrefix(line, "User-Agent: ")
				headerMap["UserAgent"] = content
			}
		} else if i > emptyLineIndex {
			body += line + "\r\n"
		}
	}

	return HTTPRequest{
		Method:    headerMap["Method"],
		Path:      headerMap["Path"],
		Version:   headerMap["Version"],
		UserAgent: headerMap["UserAgent"],
		Body:      strings.TrimSpace(body),
	}
}

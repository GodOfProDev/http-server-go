package main

import (
	"log/slog"
	"net"
	"strings"
)

func writeToConnection(header string, conn net.Conn) {
	_, err := conn.Write([]byte(header))
	if err != nil {
		slog.Error("Failed to write to the connection: ", err)
	}
}

func stripPath(str string) string {
	s := strings.Split(string(str), " ")
	content, _ := strings.CutPrefix(s[1], "/echo/")

	return content
}

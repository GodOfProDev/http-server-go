package util

import (
	"log/slog"
	"net"
)

func WriteToConnection(header string, conn net.Conn) {
	_, err := conn.Write([]byte(header))
	if err != nil {
		slog.Error("Failed to write to the connection: ", err)
	}
}

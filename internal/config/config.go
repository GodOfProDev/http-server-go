package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Directory string
	Addr      string
	Port      int
	FullAddr  string
}

func NewConfig() *Config {
	var directory string
	var addr string
	var port int
	flag.StringVar(&directory, "directory", "empty", "Sets the current directory of the server")
	flag.StringVar(&addr, "addr", "0.0.0.0", "Sets the addr of the server")
	flag.IntVar(&port, "port", 4221, "Sets the port of the server")

	flag.Parse()

	return &Config{
		Directory: directory,
		Addr:      addr,
		Port:      port,
		FullAddr:  fmt.Sprintf("%s:%d", addr, port),
	}
}

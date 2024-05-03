package config

import "flag"

type Config struct {
	Directory string
}

func NewConfig() *Config {
	var directory string
	flag.StringVar(&directory, "directory", "empty", "Set the current directory of the server")

	flag.Parse()

	return &Config{
		Directory: directory,
	}
}

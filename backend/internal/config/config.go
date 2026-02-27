package config

import "strings"

// This document loads backend configuration files like port number
type Config struct {
	Port string
}

func Load(port string) *Config {
	if port == "" {
		port = "8080" // Default port
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return &Config{
		Port: port,
	}
}

package config

import (
	"os"
	"strconv"
)

const (
	defaultPort   = 7540
	defaultWebDir = "./web"
	defaultDBFile = "scheduler.db"
)

func Port() string {
	if p := os.Getenv("TODO_PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			return p
		}
	}
	return strconv.Itoa(defaultPort)
}

func WebDir() string {
	if w := os.Getenv("TODO_WEB_DIR"); w != "" {
		return w
	}
	return defaultWebDir
}

func DBFile() string {
	if f := os.Getenv("TODO_DBFILE"); f != "" {
		return f
	}
	return defaultDBFile
}

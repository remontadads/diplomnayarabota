package server

import (
	"log"
	"net/http"
)

func Start(port, webDir string) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":" + port
	log.Printf("сервер запущен на http://localhost%s, статика из %s", addr, webDir)

	return http.ListenAndServe(addr, nil)
}

package main

import (
	"log"

	"diplomnayarabota/internal/config"
	"diplomnayarabota/internal/server"
	"diplomnayarabota/pkg/api"
	"diplomnayarabota/pkg/db"
)

func main() {
	if err := db.Init(config.DBFile()); err != nil {
		log.Fatalf("инициализация БД: %v", err)
	}
	defer db.Close()

	api.Init()

	if err := server.Start(config.Port(), config.WebDir()); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}

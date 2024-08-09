package main

import (
	"avito_bootcamp/configs"
	"avito_bootcamp/internal/app"
	"log"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}

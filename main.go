package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mrtuuro/blog-aggregator/internal/app"
	"github.com/mrtuuro/blog-aggregator/internal/config"
)

func main() {
	loadEnv()

	cfg := config.NewConfig()
	app := app.New(cfg)

	fmt.Printf("Server listening on port %s\n", cfg.Port)
	log.Fatal(app.Start())

}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

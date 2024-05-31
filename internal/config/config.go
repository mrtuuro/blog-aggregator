package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/mrtuuro/blog-aggregator/internal/database"
)

type Config struct {
	Port string
    DB *database.Queries
}

func NewConfig() *Config {
    db, err := sql.Open("postgres", os.Getenv("DB_CONN_STRING"))
    if err != nil {
        log.Fatal("Error connecting to DB!: ", err)
    }
    dbQueries := database.New(db)


	return &Config{
		Port: os.Getenv("PORT"),
        DB: dbQueries,
	}
}

package main

import (
	"database/sql"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/valverdethiago/trading-api/api"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/util"
)

var queries *db.Queries

func main() {
	config := loadConfig()
	conn := openDatabaseConnection(config)
	queries = db.New(conn)
	startServer(config, queries)
}

func loadConfig() util.Config {
	config, err := util.LoadConfig(".", "app")
	if err != nil {
		log.Fatal("Error loading application config:", err)
	}
	return config
}

func openDatabaseConnection(config util.Config) *sql.DB {
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	return conn
}

func startServer(config util.Config, querier db.Querier) {
	server := api.NewServer(queries)
	err := server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}

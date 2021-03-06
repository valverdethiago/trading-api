package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/valverdethiago/trading-api/api"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/util"
)

var queries *db.Queries

func main() {
	config, err := util.LoadConfig(".", "app")
	if err != nil {
		log.Fatal("Error loading application config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	queries = db.New(conn)

	server := api.NewServer(queries)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}

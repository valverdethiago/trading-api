package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/valverdethiago/trading-api/api"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:golang@localhost:6432/trade?sslmode=disable"
	serverAddress = "0.0.0.0:8081"
)

var queries *db.Queries

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	queries = db.New(conn)

	server := api.NewServer(queries)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}

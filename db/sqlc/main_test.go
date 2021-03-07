package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/valverdethiago/trading-api/util"
)

var testQueries Querier

func TestMain(m *testing.M) {
	config := loadConfig()
	conn := openDatabaseConnection(config)
	testQueries = New(conn)
	os.Exit(m.Run())
}

func loadConfig() util.Config {
	config, err := util.LoadConfig("../..", "test")
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

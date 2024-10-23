package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// update to env
const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:password@localhost:5432/grpc?sslmode=disable"
)

var testQueries *Queries // define as global var bc we would use it in all our unit tests
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run()) // run tests
}

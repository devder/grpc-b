package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/devder/grpc-b/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries // define as global var bc we would use it in all our unit tests
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("could not load config file: ", err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run()) // run tests
}

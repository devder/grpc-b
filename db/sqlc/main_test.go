package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/devder/grpc-b/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testStore Store // define as global var bc we would use it in all our unit tests

func TestMain(m *testing.M) {
	var err error
	_, err = util.LoadConfig("../..")
	if err != nil {
		log.Fatal("could not load config file: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), "postgresql://root:password@localhost:5432/grpc?sslmode=disable")

	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run()) // run tests
}

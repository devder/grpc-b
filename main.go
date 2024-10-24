package main

import (
	"database/sql"
	"log"

	"github.com/devder/grpc-b/api"
	db "github.com/devder/grpc-b/db/sqlc"
	_ "github.com/lib/pq"
)

// update to env
const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:password@localhost:5432/grpc?sslmode=disable"
	serverAddress  = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

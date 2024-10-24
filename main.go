package main

import (
	"database/sql"
	"log"

	"github.com/devder/grpc-b/api"
	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load config file")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

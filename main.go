package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/devder/grpc-b/api"
	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/gapi"
	"github.com/devder/grpc-b/pb"
	"github.com/devder/grpc-b/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGRPCServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGrpcAppServer(grpcServer, server)
	// consider this self documentation to help the client
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("start gRPC server at %s", lis.Addr().String())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

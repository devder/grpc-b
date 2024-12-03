package gapi

import (
	"fmt"

	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/pb"
	"github.com/devder/grpc-b/token"
	"github.com/devder/grpc-b/util"
)

type Server struct {
	pb.UnimplementedGrpcAppServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// creates new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

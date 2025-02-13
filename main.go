package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/devder/grpc-b/api"
	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/gapi"
	"github.com/devder/grpc-b/mail"
	"github.com/devder/grpc-b/pb"
	"github.com/devder/grpc-b/util"
	"github.com/devder/grpc-b/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("could not load config file")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to DB:")
	}

	if err = connPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("cannot ping DB:")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{Addr: config.RedisAddress}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store, config)
	go runGatewayServer(config, store, taskDistributor) // run in a separate go routine to avoid blocking the grpc server
	runGRPCServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create a new migrate instance:")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up: ")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config util.Config) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, config, mailer)
	log.Info().Msg("start task processor")

	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server:")
	}
}

func runGRPCServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	// interceptor
	grpcLoggerInterceptor := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLoggerInterceptor)
	pb.RegisterGrpcAppServer(grpcServer, server)
	// consider this self documentation to help the client
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	log.Info().Msgf("start gRPC server at %s", lis.Addr().String())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	// keep responses as defined in .proto (snake_case)
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterGrpcAppHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msgf("cannot register handler server: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	lis, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	log.Info().Msgf("start HTTP gateway server at %s", lis.Addr().String())

	handler := gapi.HttpLogger(mux)
	if err := http.Serve(lis, handler); err != nil {
		log.Fatal().Msgf("failed to serve HTTP gateway server: %v", err)
	}
}

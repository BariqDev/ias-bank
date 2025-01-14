package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/BariqDev/ias-bank/api"
	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/gapi"
	"github.com/BariqDev/ias-bank/mail"
	"github.com/BariqDev/ias-bank/pb"
	"github.com/BariqDev/ias-bank/util"
	"github.com/BariqDev/ias-bank/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// Load configuration file
	config, err := util.LoadConfig(".")
	ctx := context.Background()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config:")
		os.Exit(1)
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	testDbPool, err := pgxpool.New(ctx, config.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to DB:")
		os.Exit(1)
	}
	runDBMigrationUp(config.MigrationUrl, config.DBSource)

	store := db.NewStore(testDbPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisServerAddress,
	}
	taskDistributer := worker.NewRedisTaskDistributer(redisOpt)

	go runTaskProcessor(redisOpt, store, config)
	go runGrpcGatewayServer(config, store, taskDistributer)
	runGrpcServer(config, store, taskDistributer)
}

func runGrpcServer(config util.Config, store db.Store, distributer worker.TaskDistributer) {

	server, err := gapi.NewServer(config, store, distributer)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterIASBankServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener:")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server:")
	}

}

func runDBMigrationUp(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create migration:")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("Cannot migrate up:")
	}
	log.Info().Msg("Migration up success")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config util.Config) {

	mailer := mail.NewGmailSender("ias-bank", config.VerifyEmailAddress, config.VerifyEmailPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("Start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start task processor")
	}
}

func runGrpcGatewayServer(config util.Config, store db.Store, distributer worker.TaskDistributer) {
	server, err := gapi.NewServer(config, store, distributer)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

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

	err = pb.RegisterIASBankServiceHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener:")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())

	handler := gapi.HttLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gateway:")
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server:")
		os.Exit(1)
	}
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server:")
		os.Exit(1)
	}

}

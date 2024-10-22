package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	db "github.com/JorniZ/simplebank/db/sqlc"
	_ "github.com/JorniZ/simplebank/doc/statik"
	"github.com/JorniZ/simplebank/gapi"
	"github.com/JorniZ/simplebank/pb"
	"github.com/JorniZ/simplebank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config:")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to the db:")
	}

	if err := conn.Ping(); err != nil {
		log.Fatal().Err(err).Msg("cannot ping db:")
	}

	runDBMigration("file://db/migration", config.DBSource)

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating new migration: ")
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("error running db migration up: ")
	}

	log.Info().Msg("db migrated successfully")
}

func runGrpcServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating grpc listener:")
	}

	log.Info().Msgf("starting GRPC server at %s", listener.Addr().String())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("cannot start grpc server:")
	}
}

func runGatewayServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
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

	if err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server); err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create statik fs:")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating grpc listener:")
	}

	log.Info().Msgf("starting HTTP gateway server at %s", listener.Addr().String())

	handler := gapi.HttpLogger(mux)

	if err := http.Serve(listener, handler); err != nil {
		log.Fatal().Err(err).Msg("cannot start grpc server:")
	}
}

// This is old code we still wan't to keep
// func runGinServer(config util.Config, store db.Store) {
// 	server, err := api.NewServer(config, store)
// 	if err != nil {
// 		log.Fatal("cannot create server:", err.Error())
// 	}

// 	if err := server.Start(config.HTTPServerAddress); err != nil {
// 		log.Fatal("cannot start server:", err.Error())
// 	}
// }

package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/JorniZ/simplebank/api"
	db "github.com/JorniZ/simplebank/db/sqlc"
	"github.com/JorniZ/simplebank/gapi"
	"github.com/JorniZ/simplebank/pb"
	"github.com/JorniZ/simplebank/util"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db:", err.Error())
	}

	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err.Error())
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("error creating grpc listener:", err.Error())
	}

	log.Printf("starting GRPC server at %s", listener.Addr().String())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start grpc server:", err.Error())
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err.Error())
	}

	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("cannot start server:", err.Error())
	}
}

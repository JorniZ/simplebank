package main

import (
	"database/sql"
	"log"

	"github.com/JorniZ/simplebank/api"
	db "github.com/JorniZ/simplebank/db/sqlc"
	"github.com/JorniZ/simplebank/util"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
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
	server := api.NewServer(store)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err.Error())
	}
}

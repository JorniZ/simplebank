package gapi

import (
	"fmt"

	db "github.com/JorniZ/simplebank/db/sqlc"
	"github.com/JorniZ/simplebank/pb"
	"github.com/JorniZ/simplebank/token"
	"github.com/JorniZ/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(util.RandomString(32))
	if err != nil {
		return nil, fmt.Errorf("unable to create token: %s", err.Error())
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

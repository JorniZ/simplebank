package api

import (
	"fmt"
	"log"

	db "github.com/JorniZ/simplebank/db/sqlc"
	"github.com/JorniZ/simplebank/token"
	"github.com/JorniZ/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	foreignKeyViolationErrCode = "23503"
	UniqueViolationErrCode     = "23505"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			log.Fatal("error registering currency validation:", err.Error())
		}
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	accounts := router.Group("/accounts")

	accounts.POST("", server.createAccount)
	accounts.GET("/:id", server.getAccount)
	accounts.GET("", server.listAccount)
	accounts.PUT("/:id", server.updateAccount)
	accounts.DELETE("/:id", server.deleteAccount)

	transfers := router.Group("/transfers")

	transfers.POST("", server.createTransfer)

	users := router.Group("/users")

	users.POST("", server.createUser)
	users.POST("/login", server.loginUser)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

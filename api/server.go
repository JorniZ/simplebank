package api

import (
	"log"

	db "github.com/JorniZ/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	foreignKeyViolationErrCode = "23503"
	UniqueViolationErrCode     = "23505"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			log.Fatal("error registering currency validation:", err.Error())
		}
	}

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

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

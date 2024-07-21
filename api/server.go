package api

import (
	db "github.com/JorniZ/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	accounts := router.Group("/accounts")

	accounts.POST("", server.createAccount)
	accounts.GET("/:id", server.getAccount)
	accounts.GET("", server.listAccount)
	accounts.PUT("/:id", server.updateAccount)
	accounts.DELETE("/:id", server.deleteAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

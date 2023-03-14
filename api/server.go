package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/guy-ronen/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/accounts", server.createAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}


func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

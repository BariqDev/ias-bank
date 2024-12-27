package api

import (
	"github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serve http requests for application
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates new http server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}
	// routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.ListAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server

}

func (server *Server) Start(address string) error {

	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

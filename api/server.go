package api

import (
	"fmt"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/token"
	"github.com/BariqDev/ias-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serve http requests for application
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

// NewServer creates new http server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	server.setupRoutes()
	return server, nil

}

func (server *Server) setupRoutes() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/refresh", server.renewAccessToken)
	routesGroup := router.Group("/").Use(authMiddleware(server.tokenMaker))

	routesGroup.POST("/accounts", server.createAccount)
	routesGroup.GET("/accounts/:id", server.getAccount)
	routesGroup.GET("/accounts", server.ListAccounts)

	routesGroup.POST("/transfers", server.createTransfer)

	server.router = router

}

func (server *Server) Start(address string) error {

	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

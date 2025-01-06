package gapi

import (
	"fmt"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/pb"
	"github.com/BariqDev/ias-bank/token"
	"github.com/BariqDev/ias-bank/util"
)

// Server serve http requests for application
type Server struct {
	pb.UnimplementedIASBankServiceServer
	store      db.Store
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


	return server, nil

}


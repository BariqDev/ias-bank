package gapi

import (
	"fmt"
	"reflect"
	"strings"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/pb"
	"github.com/BariqDev/ias-bank/token"
	"github.com/BariqDev/ias-bank/util"
	"github.com/go-playground/validator/v10"
)

// Server serve http requests for application
type Server struct {
	pb.UnimplementedIASBankServiceServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
	validate   *validator.Validate
}

// NewServer creates new http server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	validate := validator.New()
	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
		validate:   validate,
	}

	return server, nil

}

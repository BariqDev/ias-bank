package gapi

import (
	"context"
	"errors"
	"fmt"
	"log"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/pb"
	"github.com/BariqDev/ias-bank/util"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(req.GetPassword())
    log.Printf("Received CreateUser request: %+v", req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed at hashing password: %s", err)
	}
	args := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	fmt.Println(req.GetUsername())
	fmt.Println(req.GetEmail())
	fmt.Println(req.GetFullName())

	user, err := server.store.CreateUser(ctx, args)
	if err != nil {
		fmt.Println("inside error")
		fmt.Println(err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503", "23505":
				return nil, status.Errorf(codes.AlreadyExists, "User already exists: %s", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "Failed at hashing password: %s", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}

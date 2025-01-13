package gapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/pb"
	"github.com/BariqDev/ias-bank/util"
	"github.com/BariqDev/ias-bank/worker"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := server.validateCreateUser(req)

	if violations != nil {
		badRequest := &errdetails.BadRequest{FieldViolations: violations}
		statusInvalid := status.New(codes.InvalidArgument, "Invalid Argument")
		statusDetails, err := statusInvalid.WithDetails(badRequest)
		if err != nil {
			return nil, statusInvalid.Err()
		}
		return nil, statusDetails.Err()
	}
	hashedPassword, err := util.HashPassword(req.GetPassword())
	log.Printf("Received CreateUser request: %+v", req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed at hashing password: %s", err)
	}
	args := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{

			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.distributer.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)

		},
	}

	rxResult, err := server.store.CreateUserTx(ctx, args)
	if err != nil {
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
		User: convertUser(rxResult.User),
	}

	return res, nil
}

type createUserReq struct {
	Username string `json:"username"  validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func (server *Server) validateCreateUser(req *pb.CreateUserRequest) []*errdetails.BadRequest_FieldViolation {

	user := &createUserReq{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.Email,
		FullName: req.GetFullName(),
	}
	if err := server.validate.Struct(user); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return fieldsViolations(validationErrors)
		}
		return []*errdetails.BadRequest_FieldViolation{{
			Field:       "unknown",
			Description: fmt.Sprintf("Internal validation error: %v", err),
		}}

	}

	return nil
}

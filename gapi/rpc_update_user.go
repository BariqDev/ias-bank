package gapi

import (
	"context"
	"fmt"
	"time"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/pb"
	"github.com/BariqDev/ias-bank/util"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, UnauthenticatedError(err)
	}

	violations := server.validateUpdateUser(req)
	if violations != nil {
		return nil, invalidArgument(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to update this user")
	}

	args := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{String: req.GetFullName(), Valid: req.FullName != nil},
		Email:    pgtype.Text{String: req.GetEmail(), Valid: req.Email != nil},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed at hashing password: %s", err)
		}
		args.HashedPassword = pgtype.Text{String: hashedPassword, Valid: true}
		args.PasswordChangedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	}

	user, err := server.store.UpdateUser(ctx, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed at update user: %s", err)
	}

	res := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}

type updateUserReq struct {
	Username string `json:"username"  validate:"required,alphanum"`
	Password string `json:"password" validate:"omitempty,required,min=6"`
	FullName string `json:"full_name" validate:"omitempty,required"`
	Email    string `json:"email" validate:"omitempty,required,email"`
}

func (server *Server) validateUpdateUser(req *pb.UpdateUserRequest) []*errdetails.BadRequest_FieldViolation {

	user := &updateUserReq{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
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

func invalidArgument(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "Invalid Argument")
	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}

func UnauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "Unauthenticated: %s", err)
}

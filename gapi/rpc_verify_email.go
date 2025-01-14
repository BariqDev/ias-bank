package gapi

import (
	"context"
	"fmt"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/pb"
	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {

	violations := server.validateVerifyEmail(req)

	if violations != nil {

		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, invalidArgumentError(violations)

	}
	res := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}

	return res, nil
}

type verifyEmailReq struct {
	EmailId    int64  `json:"email_id"  validate:"required,min=1"`
	SecretCode string `json:"secret_code" validate:"required,len=32"`
}

func (server *Server) validateVerifyEmail(req *pb.VerifyEmailRequest) []*errdetails.BadRequest_FieldViolation {

	user := &verifyEmailReq{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
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

package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/BariqDev/ias-bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorization     = "authorization"
	authorizationType = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}

	values := md.Get(authorization)
	if len(values) == 0 {
		return nil, fmt.Errorf("authorization token is not provided")
	}
	authHeader := values[0]

	fields := strings.Fields(authHeader)

	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	authType := strings.ToLower(fields[0])
	if authType != authorizationType {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := fields[1]

	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("access token is invalid %s", err)
	}

	return payload, nil
}

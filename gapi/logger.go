package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	starttime := time.Now()
	result, err := handler(ctx, req)

	duration := time.Since(starttime)
	statusCode := codes.Unknown

	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error()
	}
	logger.
		Str("protocol", "gRPC").
		Int("status_code", int(statusCode)).
		Str("status ", statusCode.String()).
		Dur("duration", duration).
		Msg("Received grpc request")
	return result, err
}

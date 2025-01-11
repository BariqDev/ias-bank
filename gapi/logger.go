package gapi

import (
	"context"
	"net/http"
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
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status ", statusCode.String()).
		Dur("duration", duration).
		Msg("Received grpc request")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.body = body
	return rec.ResponseWriter.Write(body)
}

func HttLogger(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		starttime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, r)
		duration := time.Since(starttime)

		logger := log.Info()

		if rec.statusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.body)
		}

		logger.
			Str("protocol", "http").
			Str("path", r.RequestURI).
			Str("method", r.Method).
			Int("status_code", int(rec.statusCode)).
			Str("status ", http.StatusText(rec.statusCode)).
			Dur("duration", duration).
			Msg("Received http request")

	})
}

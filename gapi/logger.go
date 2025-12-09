package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context, 
	req interface{}, 
	info *grpc.UnaryServerInfo, 
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)
	statusCode := codes.Unknown
	if st, err := status.FromError(err); err {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error()
	}

	logger.Str("protocol", "gRPC").
	Str("method", info.FullMethod).
	Int("status_code", int(statusCode)).
	Str("status_text", statusCode.String()).
	Dur("duration", duration).
	Msgf("received a grpc request: %s", info.FullMethod)
	return result, err
}
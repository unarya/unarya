package grpc

import (
	"context"
	"fmt"
	"net"

	"unarya/internal/shared/auth"
	"unarya/internal/shared/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryAuthInterceptor validates JWT and API key
func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if !auth.ValidateMetadata(md) {
		return nil, fmt.Errorf("unauthorized: invalid credentials")
	}
	return handler(ctx, req)
}

// UnaryLoggingInterceptor logs request info
func UnaryLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log := logging.Logger.With().Str("method", info.FullMethod).Logger()
	ctx = logging.WithContext(ctx, map[string]interface{}{"method": info.FullMethod})

	log.Info().Msg("Incoming gRPC request")
	resp, err := handler(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Request failed")
	} else {
		log.Info().Msg("Request handled successfully")
	}
	return resp, err
}

// StartGRPCServer creates a new gRPC server with interceptors
func StartGRPCServer(port string, register func(*grpc.Server)) error {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryAuthInterceptor,
			UnaryLoggingInterceptor,
		),
	)

	register(server)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	logging.Logger.Info().Msgf("gRPC server listening on :%s", port)
	return server.Serve(lis)
}

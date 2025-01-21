package helpers

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// Creates a new context embedded with auth for gRPC services
func NewServiceContext(token string) (context.Context, context.CancelFunc, error) {
	tokenGRPC, err := SignJwtForGrpc()
	if err != nil {
		return nil, nil, fmt.Errorf("error signing JWT for GRPC: %w", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	// Add the "authorization_user" metadata to the context
	md := metadata.Pairs("authorization_user", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Append additional metadata for the token (Bearer token for GRPC)
	ctxWithAuth := grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+tokenGRPC)

	return ctxWithAuth, cancel, nil
}

// Creates a new context embedded with auth for gRPC services
func NewServiceWithoutTokenContext() (context.Context, context.CancelFunc, error) {
	tokenGRPC, err := SignJwtForGrpc()
	if err != nil {
		return nil, nil, fmt.Errorf("error signing JWT for GRPC: %w", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// Append additional metadata for the token (Bearer token for GRPC)
	ctxWithAuth := grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+tokenGRPC)

	return ctxWithAuth, cancel, nil
}

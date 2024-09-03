package grpcServer

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mrkovshik/memento/internal/auth"
	"github.com/mrkovshik/memento/internal/service/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor logs details of the gRPC unary calls.
func UnaryLoggingInterceptor(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		logger.Infof("Started unary call: %s", info.FullMethod)

		// Call the handler
		resp, err := handler(ctx, req)

		// Log details
		duration := time.Since(start)
		if err != nil {
			logger.Errorf("Unary call %s finished with error: %v (duration: %v)", info.FullMethod, err, duration)
		} else {
			logger.Infof("Unary call %s finished successfully (duration: %v)", info.FullMethod, duration)
		}

		return resp, err
	}
}

// StreamLoggingInterceptor logs details of the gRPC stream calls.
func StreamLoggingInterceptor(logger *zap.SugaredLogger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		logger.Infof("Started stream call: %s", info.FullMethod)

		// Call the handler
		err := handler(srv, stream)

		// Log details
		duration := time.Since(start)
		if err != nil {
			logger.Errorf("Stream call %s finished with error: %v (duration: %v)", info.FullMethod, err, duration)
		} else {
			logger.Infof("Stream call %s finished successfully (duration: %v)", info.FullMethod, duration)
		}

		return err
	}
}

// Authenticate returns a UnaryServerInterceptor that performs authentication
func Authenticate(svc *server.BasicService, logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/api.grpc.Memento/AddUser" || info.FullMethod == "/api.grpc.Memento/GetToken" {
			return handler(ctx, req)
		}
		// Log the incoming request information
		logger.Infof("Processing request for method: %s", info.FullMethod)

		// Extract metadata from the context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Errorf("Missing metadata in context")
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		// Extract and validate the auth_token header
		values := md.Get("auth_token")
		if len(values) == 0 {
			logger.Errorf("Missing bearer token in Authorization header")
			return nil, status.Errorf(codes.Unauthenticated, "missing bearer token")
		}

		token := values[0]
		claims, err := auth.GetClaims(token, svc.Config.CryptoKey)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return nil, status.Errorf(codes.Unauthenticated, "Token expired")
			}
			logger.Errorf("Error validating token: %v", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		if claims == nil {
			logger.Errorf("Invalid token: %v", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		_, err = svc.GetUserByID(ctx, claims.UserID)
		if err != nil {
			logger.Errorf("Invalid token: %v", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		// Add claims to the context
		ctx = context.WithValue(ctx, auth.ClaimsKey, claims)

		// Log successful authentication
		logger.Infof("Successfully authenticated request for method: %s", info.FullMethod)
		return handler(ctx, req)
	}
}

// AuthenticateStream returns a StreamServerInterceptor that performs authentication
func AuthenticateStream(svc *server.BasicService, logger *zap.SugaredLogger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		// Log the incoming request information
		logger.Infof("Processing stream request for method: %s", info.FullMethod)

		// Extract metadata from the context
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			logger.Errorf("Missing metadata in context")
			return status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		// Extract and validate the auth_token header
		values := md.Get("auth_token")
		if len(values) == 0 {
			logger.Errorf("Missing bearer token in Authorization header")
			return status.Errorf(codes.Unauthenticated, "missing bearer token")
		}

		token := values[0]
		claims, err := auth.GetClaims(token, svc.Config.CryptoKey)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return status.Errorf(codes.Unauthenticated, "Token expired")
			}
			logger.Errorf("Error validating token: %v", err)
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		if claims == nil {
			logger.Errorf("Invalid token: %v", err)
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		_, err = svc.GetUserByID(stream.Context(), claims.UserID)
		if err != nil {
			logger.Errorf("Invalid token: %v", err)
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		// Add claims to the context
		newCtx := context.WithValue(stream.Context(), auth.ClaimsKey, claims)

		// Create a wrapped ServerStream with the new context
		wrappedStream := &wrappedServerStream{
			ServerStream: stream,
			ctx:          newCtx,
		}

		// Log successful authentication
		logger.Infof("Successfully authenticated stream request for method: %s", info.FullMethod)
		return handler(srv, wrappedStream)
	}
}

// wrappedServerStream wraps the ServerStream to allow modifying the context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the new context with authentication claims
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

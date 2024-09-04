package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	ClaimsKey    contextKey = "claims"
	authTokenKey string     = "auth_token"
)

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	claimsPointer, ok := ctx.Value(ClaimsKey).(*Claims)
	if !ok {
		return 0, errors.New("missing claims")
	}
	claims := *claimsPointer
	return claims.UserID, nil
}

func GetClaimsFromContext(ctx context.Context) (*Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	// Extract and validate the auth_token header
	values := md.Get(authTokenKey)
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing bearer token")
	}

	token := values[0]
	claims, err := getClaimsFromToken(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, status.Errorf(codes.Unauthenticated, "Token expired")
		}
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	if claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return claims, nil
}

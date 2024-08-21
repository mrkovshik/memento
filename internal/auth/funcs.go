package auth

import (
	"context"
	"errors"
)

const ClaimsKey string = "claims"

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	claimsPointer, ok := ctx.Value(ClaimsKey).(*Claims)
	if !ok {
		return 0, errors.New("missing claims")
	}
	claims := *claimsPointer
	return claims.UserID, nil
}

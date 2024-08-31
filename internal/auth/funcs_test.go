package auth

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testUserID1 = uint(123456)

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		want    uint
		wantErr bool
	}{
		{"1", context.WithValue(context.Background(), ClaimsKey, &Claims{
			RegisteredClaims: jwt.RegisteredClaims{},
			UserID:           testUserID1,
		}), testUserID1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserIDFromContext(tt.ctx)
			require.True(t, (err != nil) == tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

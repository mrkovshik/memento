package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	expiry  = time.Hour
	testKey = "MyAwesomeTestKey"
)

func TestBuildJWTString(t *testing.T) {
	tests := []struct {
		name   string
		userID uint
	}{
		{"1", 12346},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, errBuildJWTString := BuildJWTString(tt.userID, expiry, testKey)
			require.NoError(t, errBuildJWTString)
			claims, errGetClaims := GetClaims(token, testKey)
			require.NoError(t, errGetClaims)
			require.Equal(t, tt.userID, claims.UserID)
		})
	}
}

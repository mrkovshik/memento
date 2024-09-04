package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
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
			token, errBuildJWTString := BuildJWTString(tt.userID)
			require.NoError(t, errBuildJWTString)
			claims, errGetClaims := getClaimsFromToken(token)
			require.NoError(t, errGetClaims)
			require.Equal(t, tt.userID, claims.UserID)
		})
	}
}

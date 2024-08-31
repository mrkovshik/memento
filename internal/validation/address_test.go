package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want bool
	}{
		{"1", "localhost:8008", true},
		{"2", "localhost:80d08", false},
		{"3", "loren ipsum", false},
		{"3", "localhost-8008", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ValidateAddress(tt.addr))
		})
	}
}

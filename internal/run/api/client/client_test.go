package client_test

import (
	"errors"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		name          string
		inputErr      error
		wantSubString string
		wantWrapped   bool
	}{
		{
			name:          "Unauthenticated Error",
			inputErr:      errors.New("something Unauthenticated request"),
			wantSubString: "authentication failed",
			wantWrapped:   true,
		},
		{
			name:          "PermissionDenied Error",
			inputErr:      errors.New("access PermissionDenied"),
			wantSubString: "authentication failed",
			wantWrapped:   true,
		},
		{
			name:          "Other Error",
			inputErr:      errors.New("random error"),
			wantSubString: "random error",
			wantWrapped:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.WrapError(tt.inputErr)
			assert.Error(t, got)
			assert.Contains(t, got.Error(), tt.wantSubString)
			if tt.wantWrapped {
				assert.ErrorIs(t, got, tt.inputErr) // Checks if it wraps the original error
				assert.NotEqual(t, tt.inputErr, got) // Should not be exactly the same object (wrapped)
			} else {
				assert.Equal(t, tt.inputErr, got) // Should be exactly the same
			}
		})
	}
}

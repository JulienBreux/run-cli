package client

import (
	"fmt"
	"strings"

	"golang.org/x/oauth2/google"
)

// FindDefaultCredentials is a variable for dependency injection.
var FindDefaultCredentials = google.FindDefaultCredentials

// WrapError wraps the error with a user-friendly message if it's an authentication error.
func WrapError(err error) error {
	if strings.Contains(err.Error(), "Unauthenticated") || strings.Contains(err.Error(), "PermissionDenied") {
		return fmt.Errorf("authentication failed: %w. Tip: Ensure your 'gcloud auth application-default login' is valid and has permissions", err)
	}
	return err
}

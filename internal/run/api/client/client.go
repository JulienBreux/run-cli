/*
Copyright 2026 Julien Breux

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

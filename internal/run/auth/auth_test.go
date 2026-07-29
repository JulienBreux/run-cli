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

package auth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func TestGetInfo(t *testing.T) {
	// Create a temp directory for gcloud config
	tmpDir, err := os.MkdirTemp("", "gcloud-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create configurations directory
	configDir := filepath.Join(tmpDir, "configurations")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 1. Test with "default" config (implied when active_config is missing)
	// Create configurations/config_default
	defaultConfigContent := `
[core]
account = default@example.com
project = default-project

[run]
region = us-west1
`
	if err := os.WriteFile(filepath.Join(configDir, "config_default"), []byte(defaultConfigContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Set env var
	t.Setenv("CLOUDSDK_CONFIG", tmpDir)

	info, err := GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

	if info.User != "default@example.com" {
		t.Errorf("Expected User 'default@example.com', got '%s'", info.User)
	}
	if info.Project != "default-project" {
		t.Errorf("Expected Project 'default-project', got '%s'", info.Project)
	}
	if info.Region != "us-west1" {
		t.Errorf("Expected Region 'us-west1', got '%s'", info.Region)
	}

	// 2. Test with specific active config
	if err := os.WriteFile(filepath.Join(tmpDir, "active_config"), []byte("custom"), 0644); err != nil {
		t.Fatal(err)
	}

	customConfigContent := `
[core]
account = custom@example.com
project = custom-project
# Comment
; Another comment

[run]
region = eu-west1
`
	if err := os.WriteFile(filepath.Join(configDir, "config_custom"), []byte(customConfigContent), 0644); err != nil {
		t.Fatal(err)
	}

	info, err = GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

	if info.User != "custom@example.com" {
		t.Errorf("Expected User 'custom@example.com', got '%s'", info.User)
	}
	if info.Project != "custom-project" {
		t.Errorf("Expected Project 'custom-project', got '%s'", info.Project)
	}
	if info.Region != "eu-west1" {
		t.Errorf("Expected Region 'eu-west1', got '%s'", info.Region)
	}
}

func TestGetInfo_Defaults(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gcloud-test-defaults")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configDir := filepath.Join(tmpDir, "configurations")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Empty config, should default region
	if err := os.WriteFile(filepath.Join(configDir, "config_default"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CLOUDSDK_CONFIG", tmpDir)

	info, err := GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

	if info.Region != api_region.ALL {
		t.Errorf("Expected default Region 'all', got '%s'", info.Region)
	}
}

type mockTokenSource struct {
	token *oauth2.Token
	err   error
}

func (m *mockTokenSource) Token() (*oauth2.Token, error) {
	return m.token, m.err
}

func TestGetIDToken_CachingAndThreadSafety(t *testing.T) {
	// Save and restore the original package-level findDefaultCredentials
	origFindCreds := findDefaultCredentials
	defer func() {
		findDefaultCredentials = origFindCreds
		// Reset credentials cache
		idTokenCredsMu.Lock()
		idTokenCreds = nil
		idTokenCredsMu.Unlock()
	}()

	// Mock TokenSource that returns a token with "id_token" extra parameter
	mockTS := &mockTokenSource{
		token: (&oauth2.Token{}).WithExtra(map[string]interface{}{
			"id_token": "my-id-token",
		}),
	}

	credsDiscoveryCount := 0
	var discoveryMu sync.Mutex

	// Mock findDefaultCredentials
	findDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		discoveryMu.Lock()
		credsDiscoveryCount++
		discoveryMu.Unlock()
		return &google.Credentials{
			TokenSource: mockTS,
		}, nil
	}

	t.Run("Caching and Reuse", func(t *testing.T) {
		// Reset credentials cache
		idTokenCredsMu.Lock()
		idTokenCreds = nil
		idTokenCredsMu.Unlock()
		discoveryMu.Lock()
		credsDiscoveryCount = 0
		discoveryMu.Unlock()

		ctx := context.Background()

		// First call - should discover credentials
		token1, err := GetIDToken(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if token1 != "my-id-token" {
			t.Errorf("Expected token 'my-id-token', got '%s'", token1)
		}

		// Second call - should reuse cached credentials and not discover them again
		token2, err := GetIDToken(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if token2 != "my-id-token" {
			t.Errorf("Expected token 'my-id-token', got '%s'", token2)
		}

		// Check discovery count
		discoveryMu.Lock()
		count := credsDiscoveryCount
		discoveryMu.Unlock()
		if count != 1 {
			t.Errorf("Expected findDefaultCredentials to be called exactly once, got %d", count)
		}
	})

	t.Run("Thread Safety", func(t *testing.T) {
		// Reset credentials cache
		idTokenCredsMu.Lock()
		idTokenCreds = nil
		idTokenCredsMu.Unlock()
		discoveryMu.Lock()
		credsDiscoveryCount = 0
		discoveryMu.Unlock()

		ctx := context.Background()
		var wg sync.WaitGroup
		const numGoroutines = 10

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = GetIDToken(ctx)
			}()
		}
		wg.Wait()

		// Check discovery count - even with concurrent calls, discovery should only happen once
		discoveryMu.Lock()
		count := credsDiscoveryCount
		discoveryMu.Unlock()
		if count != 1 {
			t.Errorf("Expected findDefaultCredentials to be called exactly once concurrently, got %d", count)
		}
	})

	t.Run("Discovery Error Handling", func(t *testing.T) {
		// Reset credentials cache
		idTokenCredsMu.Lock()
		idTokenCreds = nil
		idTokenCredsMu.Unlock()

		// Mock discovery failure
		findDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("simulated discovery failure")
		}

		ctx := context.Background()
		_, err := GetIDToken(ctx)
		if err == nil {
			t.Fatal("Expected an error from discovery failure, got nil")
		}
		if !strings.Contains(err.Error(), "simulated discovery failure") {
			t.Errorf("Expected error to mention simulated discovery failure, got %v", err)
		}
	})
}

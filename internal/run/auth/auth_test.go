package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetInfo(t *testing.T) {
	// Create a temp directory for gcloud config
	tmpDir, err := os.MkdirTemp("", "gcloud-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

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
	defer os.RemoveAll(tmpDir)

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

	if info.Region != "us-central1" {
		t.Errorf("Expected default Region 'us-central1', got '%s'", info.Region)
	}
}

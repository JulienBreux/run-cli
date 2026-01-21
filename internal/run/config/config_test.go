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

package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/config"
)

func TestLoadSave(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove temporary directory: %v", err)
		}
	})

	// Set the home directory to the temporary directory
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME environment variable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("HOME"); err != nil {
			t.Fatalf("failed to unset HOME environment variable: %v", err)
		}
	})

	// 1. Test loading a non-existent config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error when loading non-existent config, but got: %v", err)
	}
	if cfg.Region != "" {
		t.Errorf("expected empty region for non-existent config, but got: %s", cfg.Region)
	}

	// 2. Test saving a config
	cfg.Region = "us-central1"
	cfg.Project = "my-project"
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// 3. Test loading an existing config
	loadedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load existing config: %v", err)
	}
	if loadedCfg.Region != "us-central1" {
		t.Errorf("expected region 'us-central1', but got: %s", loadedCfg.Region)
	}
	if loadedCfg.Project != "my-project" {
		t.Errorf("expected project 'my-project', but got: %s", loadedCfg.Project)
	}

	// 4. Test saving an updated config
	loadedCfg.Region = "europe-west1"
	loadedCfg.Project = "other-project"
	if err := loadedCfg.Save(); err != nil {
		t.Fatalf("failed to save updated config: %v", err)
	}

	// 5. Test loading the updated config
	updatedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load updated config: %v", err)
	}
	if updatedCfg.Region != "europe-west1" {
		t.Errorf("expected region 'europe-west1', but got: %s", updatedCfg.Region)
	}
	if updatedCfg.Project != "other-project" {
		t.Errorf("expected project 'other-project', but got: %s", updatedCfg.Project)
	}
}

func TestGetConfigPath(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove temporary directory: %v", err)
		}
	})

	// Set the home directory to the temporary directory
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME environment variable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("HOME"); err != nil {
			t.Fatalf("failed to unset HOME environment variable: %v", err)
		}
	})

	expectedPath := filepath.Join(tmpDir, ".run.yaml")
	actualPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	if actualPath != expectedPath {
		t.Errorf("expected config path '%s', but got: '%s'", expectedPath, actualPath)
	}
}

func TestCorruptedConfig(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove temporary directory: %v", err)
		}
	})

	// Set the home directory to the temporary directory
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME environment variable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("HOME"); err != nil {
			t.Fatalf("failed to unset HOME environment variable: %v", err)
		}
	})

	// Create a corrupted config file
	configPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	invalidYAML := []byte("region: us-central1\n  invalid-indent:")
	if err := os.WriteFile(configPath, invalidYAML, 0644); err != nil {
		t.Fatalf("failed to write corrupted config file: %v", err)
	}

	// Test loading the corrupted config
	_, err = config.Load()
	if err == nil {
		t.Fatal("expected an error when loading a corrupted config, but got nil")
	}

	if !strings.Contains(err.Error(), "cannot unmarshal config") {
		t.Errorf("expected error message to contain 'cannot unmarshal config', but got: %s", err.Error())
	}
}

func TestSaveCreatesDir(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove temporary directory: %v", err)
		}
	})

	// Set the home directory to a non-existent subdirectory of the temporary directory
	homeDir := filepath.Join(tmpDir, "nonexistent")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("failed to set HOME environment variable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("HOME"); err != nil {
			t.Fatalf("failed to unset HOME environment variable: %v", err)
		}
	})

	// Create and save a config
	cfg := &config.Config{Region: "us-east1"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Check that the config file was created
	configPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created")
	}
}

func TestGetConfigPath_Error(t *testing.T) {
	// Unset HOME to force error
	oldHome := os.Getenv("HOME")
	_ = os.Unsetenv("HOME")
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	_, err := config.GetConfigPath()
	if err == nil {
		t.Fatal("expected error when HOME is unset, but got nil")
	}
}

func TestLoad_HomeDirError(t *testing.T) {
	// Unset HOME to force error
	oldHome := os.Getenv("HOME")
	_ = os.Unsetenv("HOME")
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error when HOME is unset, but got nil")
	}
}

func TestLoad_ReadError(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpDir)
	})

	// Set HOME
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Unsetenv("HOME") }()

	configPath := filepath.Join(tmpDir, ".run.yaml")
	// Create file with no read permissions
	if err := os.WriteFile(configPath, []byte(""), 0000); err != nil {
		t.Fatalf("failed to create unreadable file: %v", err)
	}

	_, err = config.Load()
	if err == nil {
		t.Fatal("expected error when reading unreadable file, but got nil")
	}
}

func TestSave_WriteError(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpDir)
	})

	// Set HOME
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Unsetenv("HOME") }()

	// Create a directory where the config file should be
	configPath := filepath.Join(tmpDir, ".run.yaml")
	if err := os.Mkdir(configPath, 0755); err != nil {
		t.Fatalf("failed to create directory blocking config file: %v", err)
	}

	cfg := &config.Config{Region: "us-central1"}
	err = cfg.Save()
	if err == nil {
		t.Fatal("expected error when writing to a directory path, but got nil")
	}
}

func TestSave_MkdirError(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "run-cli-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(tmpDir, 0755) // Restore permissions for cleanup
		_ = os.RemoveAll(tmpDir)
	})

	// Make tmpDir read-only so we can't create subdirectories
	if err := os.Chmod(tmpDir, 0500); err != nil {
		t.Fatalf("failed to chmod temporary directory: %v", err)
	}

	// Set HOME to a subdirectory that doesn't exist
	subdir := filepath.Join(tmpDir, "subdir")
	_ = os.Setenv("HOME", subdir)
	defer func() { _ = os.Unsetenv("HOME") }()

	cfg := &config.Config{Region: "us-central1"}
	err = cfg.Save()
	if err == nil {
		t.Fatal("expected error when creating config directory in read-only parent, but got nil")
	}
}

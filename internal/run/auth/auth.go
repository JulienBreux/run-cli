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
	"bufio"
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"golang.org/x/oauth2/google"
)

var (
	scopes = []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"openid",
		"email",
	}
)

func getConfigDir() (string, error) {
	configDir := os.Getenv("CLOUDSDK_CONFIG")
	if configDir == "" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		// Check standard location for gcloud config
		configDir = filepath.Join(usr.HomeDir, ".config", "gcloud")
	}
	return configDir, nil
}

// GetInfo retrieves the current user info from gcloud config files.
func GetInfo() (info.Info, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return info.Info{}, err
	}

	// Read active config
	activeConfigPath := filepath.Join(configDir, "active_config")
	activeConfigBytes, err := os.ReadFile(activeConfigPath)
	var activeConfigName string
	if err != nil {
		// If active_config is missing, assume "default"
		activeConfigName = "default"
	} else {
		activeConfigName = strings.TrimSpace(string(activeConfigBytes))
	}

	return parseConfig(filepath.Join(configDir, "configurations", "config_"+activeConfigName))
}

func parseConfig(path string) (info.Info, error) {
	file, err := os.Open(path)
	if err != nil {
		return info.Info{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	var (
		account, project, region string
		section                  string
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = line[1 : len(line)-1]
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch section {
		case "core":
			switch key {
			case "account":
				account = val
			case "project":
				project = val
			}
		case "run":
			if key == "region" {
				region = val
			}
		}
	}

	if region == "" {
		region = api_region.ALL
	}

	return info.Info{
		User:    account,
		Project: project,
		Region:  region,
	}, nil
}

var (
	credsMu                sync.Mutex
	cachedCreds            *google.Credentials
	findDefaultCredentials = google.FindDefaultCredentials
)

// GetIDToken retrieves an identity token for the given audience using Google Cloud credentials.
// It uses a thread-safe lazy-initialization and caching mechanism via cachedCreds and credsMu.
// This avoids repeated credential discovery, filesystem config checks, and metadata server lookup latency
// (~300ms overhead) on every subsequent token retrieval, notably during concurrent proxy request forwarding.
var GetIDToken = func(ctx context.Context) (string, error) {
	credsMu.Lock()
	if cachedCreds == nil {
		// Use context.Background() during credentials discovery to ensure long-lived cached
		// credentials are not closed or invalidated due to transient/short-lived request context cancellations.
		bgCtx := context.Background()
		creds, err := findDefaultCredentials(bgCtx, scopes...)
		if err != nil {
			credsMu.Unlock()
			return "", fmt.Errorf("failed to find default credentials: %w", err)
		}
		cachedCreds = creds
	}
	creds := cachedCreds
	credsMu.Unlock()

	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", err
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("token response did not contain an id_token")
	}

	return idToken, nil
}

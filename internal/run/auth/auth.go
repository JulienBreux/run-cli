package auth

import (
	"bufio"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
)

// GetInfo retrieves the current user info from gcloud config files.
func GetInfo() (info.Info, error) {
	configDir := os.Getenv("CLOUDSDK_CONFIG")
	if configDir == "" {
		usr, err := user.Current()
		if err != nil {
			return info.Info{}, err
		}
		// Check standard location for gcloud config
		// On macOS/Linux it is usually ~/.config/gcloud
		// On Windows it is %APPDATA%/gcloud, but user.Current().HomeDir + .config is not standard for Windows.
		// However, gcloud often uses ~/.config/gcloud even on macOS.
		// Let's rely on checking ~/.config/gcloud first.
		configDir = filepath.Join(usr.HomeDir, ".config", "gcloud")
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
	defer file.Close()

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
			if key == "account" {
				account = val
			} else if key == "project" {
				project = val
			}
		case "run":
			if key == "region" {
				region = val
			}
		}
	}

	if region == "" {
		region = "us-central1"
	}

	return info.Info{
		User:    account,
		Project: project,
		Region:  region,
	}, nil
}
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

package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/auth"
)

// Info holds information about a running proxy.
type Info struct {
	Port   int
	Target string
	Server *http.Server
}

// Manager manages service proxies.
type Manager struct {
	mu      sync.Mutex
	proxies map[string]*Info
}

// NewManager creates a new proxy manager.
func NewManager() *Manager {
	return &Manager{
		proxies: make(map[string]*Info),
	}
}

// Start starts a proxy for the given service.
func (m *Manager) Start(ctx context.Context, serviceName, targetURL string) (*Info, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.proxies[serviceName]; exists {
		return nil, fmt.Errorf("proxy already running for service: %s", serviceName)
	}

	// Parse target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}

	// Find free port or let listener choose
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen on local port: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
		token, err := auth.GetIDToken(context.Background())
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+token)
		}

	}

	server := &http.Server{
		Handler: proxy,
	}

	go func() {
		_ = server.Serve(listener)
	}()

	info := &Info{
		Port:   port,
		Target: targetURL,
		Server: server,
	}
	m.proxies[serviceName] = info

	return info, nil
}

// Stop stops the proxy for the given service.
func (m *Manager) Stop(serviceName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, exists := m.proxies[serviceName]
	if !exists {
		return fmt.Errorf("no proxy running for service: %s", serviceName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := info.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop proxy server: %w", err)
	}

	delete(m.proxies, serviceName)
	return nil
}

// GetInfo returns info about a running proxy, if any.
func (m *Manager) GetInfo(serviceName string) *Info {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.proxies[serviceName]
}

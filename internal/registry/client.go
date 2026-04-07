// Package registry provides a client for the NuxtBlog plugin registry.
package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultRegistryURL = "https://raw.githubusercontent.com/nuxtblog/registry/main"

// Client communicates with the plugin registry.
type Client struct {
	baseURL string
}

// PluginMeta holds plugin metadata from the registry.
type PluginMeta struct {
	ID          string   `json:"name"`
	Title       string   `json:"title"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Runtime     string   `json:"runtime"`
	Repo        string   `json:"repo"`
	Homepage    string   `json:"homepage"`
	Icon        string   `json:"icon"`
	Tags        []string `json:"tags"`
	IsOfficial  bool     `json:"is_official"`
	License     string   `json:"license"`
	TrustLevel  string   `json:"trust_level"`
	DownloadURL string   `json:"download_url"`
}

// NewClient creates a registry client. If baseURL is empty, the default is used.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultRegistryURL
	}
	return &Client{baseURL: baseURL}
}

// ListPlugins fetches all plugins from the registry.
func (c *Client) ListPlugins() ([]PluginMeta, error) {
	url := c.baseURL + "/registry.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("registry unavailable (HTTP %d)", resp.StatusCode)
	}

	var plugins []PluginMeta
	if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return plugins, nil
}

// GetPlugin fetches plugin metadata from the registry.
func (c *Client) GetPlugin(id string) (*PluginMeta, error) {
	url := fmt.Sprintf("%s/plugins/%s/metadata.json", c.baseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("plugin %s not found (HTTP %d)", id, resp.StatusCode)
	}

	var meta PluginMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &meta, nil
}

// DownloadPlugin downloads an interpreted plugin to the target directory.
func (c *Client) DownloadPlugin(id, targetDir string) error {
	// TODO: Implement ZIP download and extraction
	return fmt.Errorf("not yet implemented: download plugin %s to %s", id, targetDir)
}

// DownloadSource downloads a compiled plugin's source to the target directory.
func (c *Client) DownloadSource(id, targetDir string) error {
	// TODO: Implement source download (git clone or ZIP)
	return fmt.Errorf("not yet implemented: download source %s to %s", id, targetDir)
}

package gist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dramxx/later/config"
)

const baseURL = "https://api.github.com"

type GistResponse struct {
	Files map[string]struct {
		Content string `json:"content"`
	} `json:"files"`
}

type PatchRequest struct {
	Files map[string]struct {
		Content string `json:"content"`
	} `json:"files"`
}

func GetInbox(cfg *config.Config) (string, error) {
	url := fmt.Sprintf("%s/gists/%s", baseURL, cfg.Gist.GistID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Gist.Token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("HTTP %d: failed to read body", resp.StatusCode)
		}
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var gistResp GistResponse
	if err := json.NewDecoder(resp.Body).Decode(&gistResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	file, ok := gistResp.Files["inbox.txt"]
	if !ok {
		return "", nil
	}

	return file.Content, nil
}

func UpdateInbox(cfg *config.Config, content string) error {
	url := fmt.Sprintf("%s/gists/%s", baseURL, cfg.Gist.GistID)

	patchReq := PatchRequest{
		Files: map[string]struct {
			Content string `json:"content"`
		}{
			"inbox.txt": {Content: content},
		},
	}

	jsonData, err := json.Marshal(patchReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Gist.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("HTTP %d: failed to read body", resp.StatusCode)
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

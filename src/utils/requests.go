package utils

import (
	"fmt"
	"io"
	"net/http"
)

func Request(url string, method string, requestBody io.Reader, authToken string, header map[string]string) ([]byte, int, error) {
	// Create request by method, url and body
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	// Bearer token authentication
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	// handle Error Message
	if !(resp.StatusCode >= 200 && resp.StatusCode < 500) {
		return nil, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, resp.StatusCode, nil
}

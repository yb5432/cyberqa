// Package openai provides an OpenAI-compatible API interface.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var SystemPrompt string = loadSystemPrompt()

// Client is an OpenAI API client.
// It is safe for concurrent use by multiple goroutines.
type Client struct {
	// config holds the client's configuration.
	config *Config

	// httpClient is the HTTP client used to send requests.
	httpClient *http.Client
}

// NewClient creates a new OpenAI API client with the given config.
// The client is safe for concurrent use by multiple goroutines.
func NewClient(config *Config) *Client {
	return &Client{
		config:     config,
		httpClient: &http.Client{},
	}
}

func loadSystemPrompt() string {
	// 使用os.ReadFile读取整个文件内容
	content, err := os.ReadFile("system_prompt.txt")
	if err != nil {
		panic(err)
	}
	return string(content)
}

// ChatCompletion sends a chat completion request to the OpenAI API.
// It takes a context for request cancellation and a ChatCompletionRequest.
// It returns a ChatCompletionResponse and an error if the request fails.
// The method automatically adds the system prompt if configured.
func (c *Client) ChatCompletion(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Add system prompt if configured
	if c.config.SystemPrompt != "" {
		systemMessage := Message{
			Role:    "system",
			Content: c.config.SystemPrompt,
		}
		request.Messages = append([]Message{systemMessage}, request.Messages...)
	}

	// Prepare the request body
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Determine the API endpoint
	apiBase := c.config.APIBase
	if apiBase == "" {
		apiBase = "https://api.openai.com/v1"
	}

	url := fmt.Sprintf("%s/chat/completions", apiBase)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	fmt.Printf("Sending request to %s with body: %s\n", url, string(body))
	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if bytes.Contains(respBody, []byte("```")) {
		respBody = bytes.ReplaceAll(respBody, []byte("```json"), []byte(""))
		respBody = bytes.ReplaceAll(respBody, []byte("```"), []byte(""))
	}
	// Parse the response
	var response ChatCompletionResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// anthropic.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go-backend/prompts"
)

const (
	AnthropicAPI = "https://api.anthropic.com/v1/messages"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicClient struct {
	apiKey string
	client *http.Client
}

type AnthropicRequest struct {
	Messages  []Message `json:"messages"`
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
}

type AnthropicResponse struct {
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (c *AnthropicClient) GetTemplateType(prompt string) (string, error) {
	req := AnthropicRequest{
		Messages: []Message{{
			Role:    "user",
			Content: prompt,
		}},
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 8000,
		System:    "Return either node or react based on what do you think this project should be. Only return a single word either 'node' or 'react'. Do not return anything extra",
	}

	resp, err := c.makeRequest(req)
	if err != nil {
		return "", err
	}

	if len(resp.Content) == 0 {
		return "", fmt.Errorf("empty response from Anthropic")
	}

	return resp.Content[0].Text, nil
}

func (c *AnthropicClient) Chat(messages []Message) (string, error) {
	respChan := make(chan *AnthropicResponse)
	errChan := make(chan error)

	go func() {
		req := AnthropicRequest{
			Messages:  messages,
			Model:     "claude-3-5-sonnet-20241022",
			MaxTokens: 8000,
			System:    prompts.GetSystemPrompt(""),
		}

		resp, err := c.makeRequest(req)
		if err != nil {
			errChan <- err
			return
		}
		respChan <- resp
	}()

	select {
	case resp := <-respChan:
		if len(resp.Content) == 0 {
			return "", fmt.Errorf("empty response from Anthropic")
		}
		return resp.Content[0].Text, nil
	case err := <-errChan:
		return "", fmt.Errorf("chat request failed: %w", err)
	}
}

func (c *AnthropicClient) makeRequest(req AnthropicRequest) (*AnthropicResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", AnthropicAPI, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var anthropicResp AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return nil, err
	}

	return &anthropicResp, nil
}

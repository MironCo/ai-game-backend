package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"rd-backend/internal/types"
)

type AIHandler struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewHandler() *AIHandler {
	return &AIHandler{
		client:  &http.Client{},
		baseURL: "https://openrouter.ai/api/v1/chat/completions",
		apiKey:  os.Getenv("OPENROUTER_API_KEY"),
	}
}

type Provider struct {
	Order          []string `json:"order,omitempty"`
	AllowFallbacks bool     `json:"allow_fallbacks,omitempty"`
}

type OpenRouterRequest struct {
	Model    string                    `json:"model"`
	Messages []types.OpenRouterMessage `json:"messages"`
	Provider *Provider                 `json:"provider,omitempty"`
}

func (h *AIHandler) GetChatCompletion(message string) (string, error) {
	request := OpenRouterRequest{
		Model: "meta-llama/llama-3.3-70b-instruct",
		Messages: []types.OpenRouterMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
		Provider: &Provider{
			Order:          []string{"OpenAI", "Together"},
			AllowFallbacks: false,
		},
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", h.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	h.addHeaders(req)

	resp, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response types.OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}

func (h *AIHandler) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("HTTP-Referer", "your-site-url")
	req.Header.Set("X-Title", "your-app-name")
}

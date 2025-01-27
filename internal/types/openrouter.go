// types/openrouter.go
package types

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterChoice struct {
	Message OpenRouterMessage `json:"message"`
}

type OpenRouterResponse struct {
	Choices []OpenRouterChoice `json:"choices"`
}

type Provider struct {
	Order          []string `json:"order,omitempty"`
	AllowFallbacks bool     `json:"allow_fallbacks,omitempty"`
}

type OpenRouterRequest struct {
	Model    string              `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
	Provider *Provider           `json:"provider,omitempty"`
}

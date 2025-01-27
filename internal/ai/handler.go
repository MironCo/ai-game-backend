package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"rd-backend/internal/ai/npc"
	"rd-backend/internal/types"
)

type AIHandler struct {
	client     *http.Client
	baseURL    string
	apiKey     string
	npcConfigs *npc.NPCs
}

func NewHandler(npcConfigs *npc.NPCs) *AIHandler {
	return &AIHandler{
		client:     &http.Client{},
		baseURL:    "https://openrouter.ai/api/v1/chat/completions",
		apiKey:     os.Getenv("OPENROUTER_API_KEY"),
		npcConfigs: npcConfigs,
	}
}

func (h *AIHandler) GetChatCompletion(message string, history []types.DBChatMessage, sender string, npcId string) (*string, error) {
	// Convert messages for the OpenRouter request
	npcPersonality := (*h.npcConfigs)[npcId]

	messages := make([]types.OpenRouterMessage, len(history)+2)

	messages[0] = types.OpenRouterMessage{
		Role:    "system",
		Content: npc.GenerateSystemPrompt(npcPersonality),
	}

	// Add history messages
	for i, msg := range history {
		role := "assistant"
		if msg.Sender == "player" {
			role = "user"
		}

		messages[i+1] = types.OpenRouterMessage{
			Role:    role,
			Content: msg.MessageText,
		}
	}

	// Add current message
	messages[len(history)+1] = types.OpenRouterMessage{
		Role:    "user",
		Content: message,
	}

	request := types.OpenRouterRequest{
		Model:    "meta-llama/llama-3.3-70b-instruct",
		Messages: messages,
		Provider: &types.Provider{
			Order:          []string{"Together", "DeepInfra"},
			AllowFallbacks: false,
		},
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", h.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	h.addHeaders(req)
	fmt.Println(req.Body)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response types.OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &response.Choices[0].Message.Content, nil
}

func (h *AIHandler) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("X-Title", "Riviera Dreams")
}

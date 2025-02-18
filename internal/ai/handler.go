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

// ModelConfig holds configuration for different model types
type ModelConfig struct {
	ModelName      string
	ProviderOrder  []string
	AllowFallbacks bool
}

var (
	LlamaConfig = ModelConfig{
		ModelName:      "meta-llama/llama-3.3-70b-instruct",
		ProviderOrder:  []string{"Together", "DeepInfra"},
		AllowFallbacks: false,
	}

	GPTConfig = ModelConfig{
		ModelName:      "openai/gpt-4-turbo",
		ProviderOrder:  []string{"OpenAI"},
		AllowFallbacks: false,
	}
)

type AIHandler struct {
	client          *http.Client
	baseURL         string
	apiKey          string
	npcConfigs      *npc.NPCs
	npcPhoneNumbers *npc.NPCNumbers
}

func NewAIHandler(npcConfigs *npc.NPCs, npcPhoneNumbers *npc.NPCNumbers) *AIHandler {
	return &AIHandler{
		client:          &http.Client{},
		baseURL:         "https://openrouter.ai/api/v1/chat/completions",
		apiKey:          os.Getenv("OPENROUTER_API_KEY"),
		npcConfigs:      npcConfigs,
		npcPhoneNumbers: npcPhoneNumbers,
	}
}

// makeOpenRouterRequest handles the common logic for making requests to OpenRouter
func (h *AIHandler) makeOpenRouterRequest(messages []types.OpenRouterMessage, modelConfig ModelConfig) (*string, error) {
	request := types.OpenRouterRequest{
		Model:    modelConfig.ModelName,
		Messages: messages,
		Provider: &types.Provider{
			Order:          modelConfig.ProviderOrder,
			AllowFallbacks: modelConfig.AllowFallbacks,
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

func (h *AIHandler) GetChatCompletion(message string, history []types.DBChatMessage, sender string, npcId string) (*string, error) {
	npcPersonality := (*h.npcConfigs)[npcId]
	messages := make([]types.OpenRouterMessage, len(history)+2)

	messages[0] = types.OpenRouterMessage{
		Role:    "system",
		Content: npc.GenerateSystemPrompt(npcPersonality),
	}

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

	messages[len(history)+1] = types.OpenRouterMessage{
		Role:    "user",
		Content: message,
	}

	return h.makeOpenRouterRequest(messages, LlamaConfig)
}

func (h *AIHandler) GetTextCompletion(message string, history []types.DBTextMessage, aiNumber string, playerNumber string) (*string, error) {
	npcPersonality := (*h.npcConfigs)[(*h.npcPhoneNumbers)[aiNumber]]
	messages := make([]types.OpenRouterMessage, len(history)+2)

	messages[0] = types.OpenRouterMessage{
		Role:    "system",
		Content: npc.GenerateSystemPrompt(npcPersonality) + ". The Player is texting you, so please respond as if you were texting with them, but keep your personality.",
	}

	for i, msg := range history {
		role := "assistant"
		if msg.SenderNumber == playerNumber {
			role = "user"
		}

		messages[i+1] = types.OpenRouterMessage{
			Role:    role,
			Content: msg.MessageText,
		}
	}

	messages[len(history)+1] = types.OpenRouterMessage{
		Role:    "user",
		Content: message,
	}

	return h.makeOpenRouterRequest(messages, LlamaConfig)
}

func (h *AIHandler) GetJSONCompletion(message string) (*string, error) {
	messages := []types.OpenRouterMessage{
		{
			Role:    "system",
			Content: "Return your response as a JSON object with no additional text or explanation: " + message,
		},
	}

	return h.makeOpenRouterRequest(messages, GPTConfig)
}

func (h *AIHandler) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("X-Title", "Riviera Dreams")
}

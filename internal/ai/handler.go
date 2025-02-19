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
		ModelName:      "openai/gpt-4o-mini", // Restored to original model
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
	if npcConfigs == nil || npcPhoneNumbers == nil {
		panic("npcConfigs and npcPhoneNumbers must not be nil")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		panic("OPENROUTER_API_KEY environment variable is not set")
	}

	return &AIHandler{
		client:          &http.Client{},
		baseURL:         "https://openrouter.ai/api/v1/chat/completions",
		apiKey:          apiKey,
		npcConfigs:      npcConfigs,
		npcPhoneNumbers: npcPhoneNumbers,
	}
}

// validateModelConfig ensures the model configuration is valid
func validateModelConfig(config ModelConfig) error {
	if config.ModelName == "" {
		return fmt.Errorf("model name cannot be empty")
	}
	if len(config.ProviderOrder) == 0 {
		return fmt.Errorf("provider order cannot be empty")
	}
	return nil
}

// makeOpenRouterRequest handles the common logic for making requests to OpenRouter
func (h *AIHandler) makeOpenRouterRequest(messages []types.OpenRouterMessage, modelConfig ModelConfig) (*string, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("messages array cannot be empty")
	}

	if err := validateModelConfig(modelConfig); err != nil {
		return nil, fmt.Errorf("invalid model configuration: %w", err)
	}

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
	//log.Println(req.Body)

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

func (h *AIHandler) GetChatCompletion(message string, history []types.DBChatMessage, eventHistory []types.DBPlayerEvent, sender string, npcId string) (*string, error) {
	if message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	npcPersonality, exists := (*h.npcConfigs)[npcId]
	if !exists {
		return nil, fmt.Errorf("NPC with ID %s not found", npcId)
	}

	// Initialize with capacity for system message + history + current message
	messages := make([]types.OpenRouterMessage, 0, len(history)+2)

	messages = append(messages, types.OpenRouterMessage{
		Role:    "system",
		Content: npc.GenerateSystemPromptWithEvents(npcPersonality, eventHistory),
	})

	// Add history messages if present
	if len(history) > 0 {
		for _, msg := range history {
			role := "assistant"
			if msg.Sender == "player" {
				role = "user"
			}

			messages = append(messages, types.OpenRouterMessage{
				Role:    role,
				Content: msg.MessageText,
			})
		}
	}

	// Add current message
	messages = append(messages, types.OpenRouterMessage{
		Role:    "user",
		Content: message,
	})

	return h.makeOpenRouterRequest(messages, LlamaConfig)
}

func (h *AIHandler) GetTextCompletion(message string, history []types.DBTextMessage, aiNumber string, playerNumber string) (*string, error) {
	if message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	npcId, exists := (*h.npcPhoneNumbers)[aiNumber]
	if !exists {
		return nil, fmt.Errorf("no NPC found for number %s", aiNumber)
	}

	npcPersonality, exists := (*h.npcConfigs)[npcId]
	if !exists {
		return nil, fmt.Errorf("NPC with ID %s not found", npcId)
	}

	// Initialize with capacity for system message + history + current message
	messages := make([]types.OpenRouterMessage, 0, len(history)+2)

	messages = append(messages, types.OpenRouterMessage{
		Role:    "system",
		Content: npc.GenerateSystemPrompt(npcPersonality) + ". The Player is texting you, so please respond as if you were texting with them, but keep your personality.",
	})

	// Add history messages if present
	if len(history) > 0 {
		for _, msg := range history {
			role := "assistant"
			if msg.SenderNumber == playerNumber {
				role = "user"
			}

			messages = append(messages, types.OpenRouterMessage{
				Role:    role,
				Content: msg.MessageText,
			})
		}
	}

	// Add current message
	messages = append(messages, types.OpenRouterMessage{
		Role:    "user",
		Content: message,
	})

	return h.makeOpenRouterRequest(messages, LlamaConfig)
}

func (h *AIHandler) GetJSONCompletion(message string) (*string, error) {
	if message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	messages := []types.OpenRouterMessage{
		{
			Role:    "system",
			Content: "Return your response as a JSON object with no additional text or explanation: " + message,
		},
	}

	return h.makeOpenRouterRequest(messages, GPTConfig)
}

func (h *AIHandler) GetDescriptionCompletion(message string) (*string, error) {
	if message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	messages := []types.OpenRouterMessage{
		{
			Role:    "system",
			Content: "Please summarize the info in this JSON object as a short sentence, describing what the player did. E.G: The player.... " + message,
		},
	}

	return h.makeOpenRouterRequest(messages, GPTConfig)
}

func (h *AIHandler) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("X-Title", "Riviera Dreams")
}

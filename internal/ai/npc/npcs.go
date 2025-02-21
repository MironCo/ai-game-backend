package npc

import (
	"encoding/json"
	"fmt"
	"os"
	"rd-backend/internal/types"
	"strings"
)

type NPCs map[string]types.NPC
type NPCNumbers map[string]string

func LoadNPCConfig(path string) (NPCs, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var npcs NPCs
	err = json.Unmarshal(data, &npcs)
	return npcs, err
}

func BuildPhoneIndex(npcs map[string]types.NPC) NPCNumbers {
	index := make(NPCNumbers)
	for _, npc := range npcs {
		index[npc.PhoneNumber] = npc.ID
	}
	return index
}

func GenerateSystemPrompt(npc types.NPC) string {
	return fmt.Sprintf(
		"You're %s! You're working on %s in %s. Quick bio: %s "+
			"Your friends would describe you as %s. "+
			"People can't help but notice how you %s. "+
			"These days, you're focused on %s. "+
			"When chatting, %s. "+
			"Remember to be natural and let your personality shine - no need to stick to formal speech patterns!",
		npc.Name,
		npc.Occupation,
		npc.Location,
		npc.Backstory,
		strings.Join(npc.Traits, ", "),
		strings.Join(npc.Quirks, " and "),
		npc.Goals,
		npc.SpeechStyle,
	)
}

func GenerateSystemPromptWithEvents(npc types.NPC, events []types.DBPlayerEvent) string {
	// Build the base prompt
	basePrompt := fmt.Sprintf(
		"You're %s! You're working on %s in %s. Quick bio: %s "+
			"Your friends would describe you as %s. "+
			"People can't help but notice how you %s. "+
			"These days, you're focused on %s. "+
			"When chatting, %s.",
		npc.Name,
		npc.Occupation,
		npc.Location,
		npc.Backstory,
		strings.Join(npc.Traits, ", "),
		strings.Join(npc.Quirks, " and "),
		npc.Goals,
		npc.SpeechStyle,
	)

	// If there are events, add them to the prompt
	if len(events) > 0 {
		eventDetails := make([]string, 0, len(events))
		for _, event := range events {
			eventDetails = append(eventDetails, event.EventDetails)
		}

		basePrompt += fmt.Sprintf(
			"\nThese are the things that the player has done recently, use these to inform your response: %s",
			strings.Join(eventDetails, "; "),
		)
	}

	// Add the personality reminder
	basePrompt += "\nRemember to be natural and let your personality shine - no need to stick to formal speech patterns!"

	return basePrompt
}

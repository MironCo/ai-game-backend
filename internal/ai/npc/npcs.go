package npc

import (
	"encoding/json"
	"fmt"
	"os"
	"rd-backend/internal/types"
	"strings"
)

type NPCs map[string]types.NPC

func LoadNPCConfig(path string) (NPCs, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var npcs NPCs
	err = json.Unmarshal(data, &npcs)
	return npcs, err
}

func GenerateSystemPrompt(npc types.NPC) string {
	return fmt.Sprintf(
		"You are %s, a %s in %s. %s. Your personality is %s, and you're known for %s. "+
			"Your ultimate goal is to %s. When speaking, you %s.",
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

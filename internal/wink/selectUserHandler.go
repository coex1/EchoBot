package wink

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	selectedUsersMap = make(map[string][]string)
	MinValues        int
	MaxValues        int
)


func HandleSelectMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Map 변수
	selectedUsersMap[i.GuildID] = i.MessageComponentData().Values

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 상호작용 지연
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to select menu interaction:", err)
	}
}

package mafia

import (
	// system packages
	"log"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external package
	dgo "github.com/bwmarrin/discordgo"
)

const MAFIA_MIN_LIST_CNT = 3
const MAX_MEMBER_GET int = 50
const QUERY_STRING string = ""

func SelectMenu(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	// Map 변수
	// get currently selected users, and put values to selectedUsersMap
	guild.Mafia.SelectedUsersMap[event.GuildID] = event.MessageComponentData().Values

	err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
		// 상호작용 지연
		Type: dgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to select menu interaction:", err)
	}
}

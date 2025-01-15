package mafia

import (
	// system packages
	"log"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external package
	dgo "github.com/bwmarrin/discordgo"
)

// Start
var start_selectMenu dgo.SelectMenu = dgo.SelectMenu{
	CustomID:    "mafia_Start_listUpdate",
	Placeholder: "사용자를 선택해 주세요!",
}

// on interaction event 'mafia_Start_listUpdate'
func Start_listUpdate(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.SelectedUsers[i.GuildID] = i.MessageComponentData().Values

	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		// 상호작용 지연
		Type: dgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to select menu interaction:", err)
	}
}

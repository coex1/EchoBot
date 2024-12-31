package wink

import (
	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

// on interaction event 'wink_Game_listUpdate'
func Game_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  guild.Wink.UserSelection[event.User.GlobalName] = event.MessageComponentData().Values[0]
}

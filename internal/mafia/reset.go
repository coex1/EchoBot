package mafia

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
)

func Reset(guild *data.Guild) []dgo.SelectMenuOption {
	var AliveUserID []dgo.SelectMenuOption
	guild.Mafia.VoteMap = make(map[string]string)
	guild.Mafia.VoteCount = make(map[string]int)

	// Init alive user SelectMenuOption
	for _, player := range guild.Mafia.Players {
		if player.IsAlive {
			AliveUserID = append(AliveUserID, dgo.SelectMenuOption{
				Label: player.GlobalName,
				Value: player.ID,
			})
		}
	}
	return AliveUserID
}

package mafia

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
)

func Reset(guild *data.Guild) []dgo.SelectMenuOption {
	var AliveUserID []dgo.SelectMenuOption
	guild.Mafia.TempVoteMap = make(map[string]string)
	guild.Mafia.VoteMap = make(map[string]string)
	guild.Mafia.VoteCount = make(map[string]int)

	guild.Mafia.MafiaTargetMap = make(map[string]string)

	guild.Mafia.MafiaTarget = ""
	guild.Mafia.PoliceTarget = ""
	guild.Mafia.DoctorTarget = ""

	guild.Mafia.NightActionDone = make(map[string]bool)
	guild.Mafia.NightActionDone["Mafia"] = false
	guild.Mafia.NightActionDone["Police"] = false
	guild.Mafia.NightActionDone["Doctor"] = false

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

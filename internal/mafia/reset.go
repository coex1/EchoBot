package mafia

import (
	"math/rand"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
)

func Reset(guild *data.Guild) []dgo.SelectMenuOption {
	var AliveUserID []dgo.SelectMenuOption
	guild.Mafia.TempVoteMap = make(map[string]string)
	guild.Mafia.VoteMap = make(map[string]string)
	guild.Mafia.VoteCount = make(map[string]int)
	guild.Mafia.VoteSkip = []string{}

	guild.Mafia.MafiaTargetMap = make(map[string]string)

	guild.Mafia.MafiaTarget = ""
	guild.Mafia.PoliceTarget = ""
	guild.Mafia.DoctorTarget = ""

	guild.Mafia.NightActionDone = make(map[string]bool)
	guild.Mafia.NightActionDone["Mafia"] = false
	guild.Mafia.NightActionDone["Police"] = false
	guild.Mafia.NightActionDone["Doctor"] = false
	guild.Mafia.NightActionDone["Citizen"] = false

	guild.Mafia.CitizenPhrases = make(map[string]string)
	guild.Mafia.CitizenReady = make(map[string]bool)

	// Init alive user SelectMenuOption
	for _, p := range guild.Mafia.Players {
		if p.IsAlive {
			if p.Role == "Citizen" {
				guild.Mafia.CitizenPhrases[p.ID] = guild.Mafia.SleepPhrases[rand.Intn(len(guild.Mafia.SleepPhrases))]
			}
			AliveUserID = append(AliveUserID, dgo.SelectMenuOption{
				Label: p.GlobalName,
				Value: p.ID,
			})
		}
	}
	return AliveUserID
}

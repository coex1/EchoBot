package mafia

import (
	// system packages
	"log"
	"math/rand"
	"time"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

func StartButton(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	mafiaCount := event.ApplicationCommandData().Options[0].IntValue()
	policeCount := event.ApplicationCommandData().Options[1].IntValue()
	doctorCount := event.ApplicationCommandData().Options[2].IntValue()

	var mafiaSelected, policeSelected, doctorSelected []string
	totalCount := int(mafiaCount + policeCount + doctorCount)

	// tempSelectedMembers := selectedUsersMap[event.GuildID]
	if len(guild.Mafia.SelectedUsersMap) == 0 {
		log.Println("No members selected.")
		return
	}

	if totalCount > len(guild.Mafia.SelectedUsersMap) {
		log.Println("Exceeded number")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]string, 0)
	copy(shuffled, guild.Mafia.SelectedUsersMap[event.GuildID])
	r.Shuffle(len(shuffled), func(event, j int) {
		shuffled[event], shuffled[j] = shuffled[j], shuffled[event]
	})
	mafiaSelected = shuffled[:mafiaCount]
	shuffled = shuffled[mafiaCount:]
	policeSelected = shuffled[:policeCount]
	shuffled = shuffled[policeCount:]
	doctorSelected = shuffled[:doctorCount]

	var message string
	for _, id := range guild.Mafia.SelectedUsersMap[event.GuildID] {
		if general.Contains(mafiaSelected, id) {
			message = "당신은 마피아 입니다!"
		} else if general.Contains(policeSelected, id) {
			message = "당신은 경찰 입니다!"
		} else if general.Contains(doctorSelected, id) {
			message = "당신은 의사 입니다!"
		} else {
			message = "당신은 시민 입니다!"
		}
		general.SendDM(s, id, message)
	}
}

package mafia

// system packages
import (
	"log"
	"math/rand"
	"time"
)

// internal imports
import (
	"github.com/coex1/EchoBot/internal/discord"
)

// external packages
import (
	dgo "github.com/bwmarrin/discordgo"
)

func StartButton(s *dgo.Session, event *dgo.InteractionCreate) {
	mafiaCount := event.ApplicationCommandData().Options[0].IntValue()
	policeCount := event.ApplicationCommandData().Options[1].IntValue()
	doctorCount := event.ApplicationCommandData().Options[2].IntValue()

	var mafiaSelected, policeSelected, doctorSelected []string
	totalCount := int(mafiaCount + policeCount + doctorCount)

	tempSelectedMembers := selectedUsersMap[event.GuildID]
	if len(tempSelectedMembers) == 0 {
		log.Println("No members selected.")
		return
	}

	if totalCount > len(tempSelectedMembers) {
		log.Println("Exceeded number")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]string, 0)
	copy(shuffled, tempSelectedMembers)
	r.Shuffle(len(shuffled), func(event, j int) {
		shuffled[event], shuffled[j] = shuffled[j], shuffled[event]
	})
	mafiaSelected = shuffled[:mafiaCount]
	shuffled = shuffled[mafiaCount:]
	policeSelected = shuffled[:policeCount]
	shuffled = shuffled[policeCount:]
	doctorSelected = shuffled[:doctorCount]

	var message string
	for _, id := range tempSelectedMembers {
		if contains(mafiaSelected, id) {
			message = "당신은 마피아 입니다!"
		} else if contains(policeSelected, id) {
			message = "당신은 경찰 입니다!"
		} else if contains(doctorSelected, id) {
			message = "당신은 의사 입니다!"
		} else {
			message = "당신은 시민 입니다!"
		}
		discord.SendDM(s, id, message)
	}
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
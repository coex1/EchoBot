package discord

// system packages
import (
	"log"
	"math/rand"
	"time"
)

// external package
import (
	dgo "github.com/bwmarrin/discordgo"
)

var (
  // 윙크 받아서 버튼은 클릭 한 사용자들 
	MafiaSelectedUsersMap = make(map[string][]string)
)

func Mafia_HandleSelectMenu(s *dgo.Session, event *dgo.InteractionCreate) {
	// Map 변수
  // get currently selected users, and put values to selectedUsersMap
	MafiaSelectedUsersMap[event.GuildID] = event.MessageComponentData().Values

	err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
		// 상호작용 지연
		Type: dgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to select menu interaction:", err)
	}
}

func Mafia_HandleStartButton(s *dgo.Session, event *dgo.InteractionCreate) {
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
		SendDM(s, id, message)
	}
}

package mafia

import (
	"math/rand"
	"time"

	// internal packages
	"github.com/coex1/EchoBot/internal/general"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

func assignRole(s *dgo.Session, players []string, mafiaCount int, policeCount int, doctorCount int) (mafiaID []string, policeID []string, doctorID []string, citizenID []string) {
	shuffled := make([]string, len(players))
	copy(shuffled, players)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	mafiaID = shuffled[:mafiaCount]
	shuffled = shuffled[mafiaCount:]
	policeID = shuffled[:policeCount]
	shuffled = shuffled[policeCount:]
	citizenID = shuffled[:doctorCount]

	var message string
	for _, id := range players {
		if general.Contains(mafiaID, id) {
			message = "당신은 마피아 입니다!"
		} else if general.Contains(policeID, id) {
			message = "당신은 경찰 입니다!"
		} else if general.Contains(citizenID, id) {
			message = "당신은 의사 입니다!"
		} else {
			message = "당신은 시민 입니다!"
		}
		general.SendDM(s, id, message)
	}
	return
}

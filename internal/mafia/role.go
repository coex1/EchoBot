package mafia

import (
	"log"
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/general"
)

// 플레이어 역할 배정
func assignRole(players []string, mafiaCount int, policeCount int, doctorCount int) (mafiaID []string, policeID []string, doctorID []string, citizenID []string) {
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

	return
}

// 역할별 DM 발송
func sendRoleDMs(s *dgo.Session, players, mafiaIDs, policeIDs, doctorIDs []string) {
	for _, id := range players {
		var message string

		switch {
		case general.Contains(mafiaIDs, id):
			message = "당신은 **마피아**입니다! 밤마다 시민을 처치할 수 있습니다."
		case general.Contains(policeIDs, id):
			message = "당신은 **경찰**입니다! 밤마다 한 명의 신원을 확인할 수 있습니다."
		case general.Contains(doctorIDs, id):
			message = "당신은 **의사**입니다! 밤마다 한 명을 보호할 수 있습니다."
		default:
			message = "당신은 **시민**입니다! 마피아를 찾아내세요."
		}

		// 개인 DM 발송
		err := general.SendDM(s, id, message)
		if err != nil {
			log.Printf("Failed to send DM to user %s: %v\n", id, err)
		}
	}
}

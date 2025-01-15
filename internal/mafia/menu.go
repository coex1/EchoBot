package mafia

import (
	// system packages
	"log"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external package
	dgo "github.com/bwmarrin/discordgo"
)

// on interaction event 'mafia_Start_listUpdate'
// 드롭다운 선택 시 실행
func Start_listUpdate(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {

	// 선택한 플레이어 저장
	guild.Mafia.SelectedUsers[i.GuildID] = i.MessageComponentData().Values

	// 선택한 플레이어 목록 출력
	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseUpdateMessage,
	})
	if err != nil {
		log.Println("플레이어 목록 출력 실패:", err)
	}
}

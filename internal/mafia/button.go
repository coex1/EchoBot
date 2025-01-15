package mafia

import (
	// system packages

	// internal packages
	"log"

	"github.com/coex1/EchoBot/internal/data"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

// on interaction event 'mafia_Start_Button'
func Start_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {

	players := guild.Mafia.SelectedUsers[i.GuildID]

	var mafiaCount = int(i.ApplicationCommandData().Options[0].IntValue())
	var policeCount = int(i.ApplicationCommandData().Options[1].IntValue())
	var doctorCount = int(i.ApplicationCommandData().Options[2].IntValue())

	if len(players) < MIN_PLAYER_CNT {
		log.Println("Invalid player count!")
		return
	}
	if mafiaCount+policeCount+doctorCount > len(players) {
		log.Println("Exceeded count")
		return
	}

	mafiaIDs, policeIDs, doctorIDs, citizenIDs := assignRole(players, mafiaCount, policeCount, doctorCount)

	guild.Mafia.MafiaList = mafiaIDs
	guild.Mafia.PoliceList = policeIDs
	guild.Mafia.DoctorList = doctorIDs
	guild.Mafia.CitizenList = citizenIDs

	sendRoleDMs(s, players, mafiaIDs, policeIDs, doctorIDs)

	// 게임 시작 메시지 전송
	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "마피아 게임이 시작되었습니다! 역할이 개별 DM으로 전송되었습니다.",
		},
	})
	if err != nil {
		log.Println("Error sending game start message:", err)
	}
}

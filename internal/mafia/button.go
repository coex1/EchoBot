package mafia

import (
	// system packages

	// internal packages
	"log"

	"github.com/coex1/EchoBot/internal/data"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

// start_button
var start_buttonRow dgo.ActionsRow = dgo.ActionsRow{
	Components: []dgo.MessageComponent{
		&dgo.Button{
			Label:    "게임 시작",              // 버튼 텍스트
			Style:    dgo.PrimaryButton,    // 버튼 스타일
			CustomID: "mafia_Start_Button", // 버튼 클릭 시 처리할 ID
		},
	},
}

// on interaction event 'mafia_Start_Button'
func Start_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	var mafiaCount = int(i.ApplicationCommandData().Options[0].IntValue())
	var policeCount = int(i.ApplicationCommandData().Options[1].IntValue())
	var doctorCount = int(i.ApplicationCommandData().Options[2].IntValue())

	players := guild.Mafia.SelectedUsers[i.GuildID]
	if len(players) < MIN_PLAYER_CNT {
		log.Println("Invalid player count!")
		return
	}
	if mafiaCount+policeCount+doctorCount > len(players) {
		log.Println("Exceeded count")
	}

	guild.Mafia.MafiaList,
		guild.Mafia.PoliceList,
		guild.Mafia.DoctorList,
		guild.Mafia.CitizenList = assignRole(s, guild.Mafia.SelectedUsers[i.GuildID], mafiaCount, policeCount, doctorCount)
}

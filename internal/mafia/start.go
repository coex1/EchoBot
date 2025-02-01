package mafia

import (
	// system packages
	"log"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external package
	dgo "github.com/bwmarrin/discordgo"
)

// 드롭다운 선택 시 실행
func Start_listUpdate(i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.SelectedUsersID = i.MessageComponentData().Values
}

// on interaction event 'mafia_Start_Button'
func Start_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	// Init
	guild.Mafia.State = true
	guild.Mafia.Day = 1
	guild.Mafia.VoteMap = make(map[string]string)
	guild.Mafia.VoteCount = make(map[string]int)
	guild.Mafia.ReadyMap = make(map[string]bool)
	guild.Mafia.Timer = 0 // TODO : function

	for _, id := range guild.Mafia.SelectedUsersID {
		guild.Mafia.ReadyMap[id] = false
	}

	var numMafia = guild.Mafia.NumMafia
	var numPolice = guild.Mafia.NumPolice
	var numDoctor = guild.Mafia.NumDoctor

	if len(guild.Mafia.SelectedUsersID) < MIN_PLAYER_CNT {
		log.Println("Invalid player count!")
		return
	}
	if numMafia+numPolice+numDoctor > len(guild.Mafia.SelectedUsersID) {
		log.Println("Exceeded count")
		return
	}

	mafiaIDs, policeIDs, doctorIDs, citizenIDs :=
		sendStartMessage(s, guild.Mafia.SelectedUsersID, numMafia, numPolice, numDoctor)

	guild.Mafia.MafiaList = mafiaIDs
	guild.Mafia.PoliceList = policeIDs
	guild.Mafia.DoctorList = doctorIDs
	guild.Mafia.CitizenList = citizenIDs

	//
	Role_Message(s, i, guild)
}

func Rdy_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	checkAllPlayersReady := func(readyMap map[string]bool) bool {
		for _, ready := range readyMap {
			if !ready {
				return false
			}
		}
		return true
	}
	guild.Mafia.ReadyMap[i.User.ID] = true

	if checkAllPlayersReady(guild.Mafia.ReadyMap) { // 게임 시작
		guild.Mafia.State = true                               // 아침 설정
		guild.Mafia.AliveUsersID = guild.Mafia.SelectedUsersID // 생존 유저
		Day_Message(s, i, guild)
	} else {
		err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			log.Printf("Failed to delayed respone: %v\n", err)
			return
		}
	}
}

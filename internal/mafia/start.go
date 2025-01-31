package mafia

import (
	// system packages
	"log"
	"strconv"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"

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
	guild.Mafia.UserDMChannels = make(map[string]string)
	guild.Mafia.VoteMap = make(map[string]string)
	guild.Mafia.VoteCount = make(map[string]int)

	players := guild.Mafia.SelectedUsersID

	for _, id := range players {
		guild.Mafia.UserDMChannels[id] = ""
	}

	guild.Mafia.AliveUsersID = guild.Mafia.SelectedUsersID

	var numMafia = guild.Mafia.NumMafia
	var numPolice = guild.Mafia.NumPolice
	var NumDoctor = guild.Mafia.NumDoctor

	if len(players) < MIN_PLAYER_CNT {
		log.Println("Invalid player count!")
		return
	}
	if numMafia+numPolice+NumDoctor > len(players) {
		log.Println("Exceeded count")
		return
	}

	mafiaIDs, policeIDs, doctorIDs, citizenIDs :=
		sendStartMessage(s, guild, players, numMafia, numPolice, NumDoctor)

	guild.Mafia.MafiaList = mafiaIDs
	guild.Mafia.PoliceList = policeIDs
	guild.Mafia.DoctorList = doctorIDs
	guild.Mafia.CitizenList = citizenIDs

	//
	Role_Message(s, i, guild)
}

func Rdy_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.ReadyID = append(guild.Mafia.ReadyID, i.User.ID)
	numReady := strconv.Itoa(len(guild.Mafia.ReadyID))
	numTotal := strconv.Itoa(len(guild.Mafia.AliveUsersID))

	if numReady == numTotal { // 게임 시작
		guild.Mafia.State = true // 아침 설정
		Start_Message(s, i, guild)
	} else {
		msg := "준비 완료!" + numReady + " / " + numTotal
		err := general.SendDM(s, i.User.ID, msg)
		if err != nil {
			log.Printf("Failed to send DM to user %s: %v\n", i.User.ID, err)
		}
	}
}

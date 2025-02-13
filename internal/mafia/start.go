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
	guild.Mafia.Players = make(map[string]*data.MafiaPlayer)
	guild.Mafia.State = true
	guild.Mafia.Day = 1
	guild.Mafia.VoteMap = make(map[string]string)
	guild.Mafia.VoteCount = make(map[string]int)
	guild.Mafia.ReadyMap = make(map[string]bool)
	guild.Mafia.Timer = 0 // TODO : function

	for _, id := range guild.Mafia.SelectedUsersID {
		member, err := s.GuildMember(i.GuildID, id)
		if err != nil {
			log.Fatalf("Failed getting members [%v]", err)
			return
		}
		guild.Mafia.Players[id] = &data.MafiaPlayer{
			ID:         id,
			GlobalName: member.User.GlobalName,
			IsAlive:    true,
		}

		guild.Mafia.ReadyMap[id] = false
		guild.Mafia.AliveUsersID = append(guild.Mafia.AliveUsersID, member.User.ID)
		guild.Mafia.AliveUsersIDMap = append(guild.Mafia.AliveUsersIDMap, dgo.SelectMenuOption{
			Label: member.User.GlobalName,
			Value: member.User.ID,
		})
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
		sendStartMessage(s, guild, guild.Mafia.SelectedUsersID, numMafia, numPolice, numDoctor)

	guild.Mafia.MafiaList = mafiaIDs
	guild.Mafia.PoliceList = policeIDs
	guild.Mafia.DoctorList = doctorIDs
	guild.Mafia.CitizenList = citizenIDs

	// 역할 개별 전송
	Role_Message(s, i, guild)
}

func Ready_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	allPlayersReady := func(players []string) bool {
		for _, id := range players {
			if !guild.Mafia.ReadyMap[id] { // 한 명이라도 Ready가 아니면 false 반환
				return false
			}
		}
		return true
	}

	guild.Mafia.ReadyMap[i.User.ID] = true

	log.Printf("User %s is ready!", i.User.ID)
	// 모든 유저가 준비 완료되었는지 확인 후 게임 시작
	if allPlayersReady(guild.Mafia.SelectedUsersID) {
		Day_Message(s, guild)
	}
}

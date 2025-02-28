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
	guild.Mafia.Day = 1
	guild.Mafia.State = true
	guild.Mafia.ReadyMap = make(map[string]bool)
	guild.Mafia.Timer = 180 // TODO : function

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
	}

	// 투표 필요 정보
	Reset(guild)

	if len(guild.Mafia.SelectedUsersID) < MIN_PLAYER_CNT {
		log.Println("Invalid player count!")
		return
	}
	if guild.Mafia.NumMafia+guild.Mafia.NumPolice+guild.Mafia.NumDoctor > len(guild.Mafia.SelectedUsersID) {
		log.Println("Exceeded count")
		return
	}

	Game_Process(s, i, guild)
}

func Start_Message(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	_, err := s.ChannelMessageSendComplex(i.ChannelID, &dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			{
				Title:       "역할이 배정되었습니다!",
				Description: "**역할과 진행은 개별 DM**을 확인해주세요.",
				Color:       0xFFFFFF,
			},
		},
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					&dgo.Button{
						Label:    "게임 재시작",
						Style:    dgo.SuccessButton,
						CustomID: "mafia_restart",
					},
					&dgo.Button{
						Label:    "게임 종료",
						Style:    dgo.DangerButton,
						CustomID: "mafia_end",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to send DM to users [%v]", err)
	}
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

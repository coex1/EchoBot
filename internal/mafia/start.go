package mafia

import (
	// system packages
	"log"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external package
	dgo "github.com/bwmarrin/discordgo"
)

// 드롭다운 메뉴 선택 시
func Start_listUpdate(i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.SelectedUsersID = i.MessageComponentData().Values
}

// on interaction event 'mafia_Start_Button'
func Start_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	// 세션 변수 초기화
	guild.Mafia.Players = make(map[string]*data.MafiaPlayer)
	guild.Mafia.Day = 1
	guild.Mafia.State = true
	guild.Mafia.ReadyMap = make(map[string]bool)
	guild.Mafia.Timer = 600 // TODO : function

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

	if len(guild.Mafia.SelectedUsersID) < MIN_PLAYER_CNT {
		log.Println("Invalid player count!")
		return
	}
	if guild.Mafia.NumMafia+guild.Mafia.NumPolice+guild.Mafia.NumDoctor > len(guild.Mafia.SelectedUsersID) {
		log.Println("Exceeded count")
		return
	}

	// 게임 설명 전달 (채널)
	Start_Message(s, i, guild)

	// 역할 공지 (개별)
	Role_Message(s, guild)

	// 아침 시작
	Day_Message(s, guild)
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

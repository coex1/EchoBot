package mafia

import (
	// system packages
	"log"
	"math/rand"
	"time"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"

	// external package
	dgo "github.com/bwmarrin/discordgo"
)

// 드롭다운 선택 시 실행
func Start_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.SelectedUsersID = event.MessageComponentData().Values
}

// on interaction event 'mafia_Start_Button'
func Start_Button(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	players := guild.Mafia.SelectedUsersID

	if len(players) < MIN_PLAYER_CNT {
		log.Println("Invalid player count!")
		return
	}
	var numMafia = guild.Mafia.NumMafia
	var numPolice = guild.Mafia.NumPolice
	var NumDoctor = guild.Mafia.NumDoctor

	if numMafia+numPolice+NumDoctor > len(players) {
		log.Println("Exceeded count")
		return
	}
	mafiaIDs, policeIDs, doctorIDs, citizenIDs := sendPlayersStartMessage(s, players, numMafia, numPolice, NumDoctor)

	guild.Mafia.MafiaList = mafiaIDs
	guild.Mafia.PoliceList = policeIDs
	guild.Mafia.DoctorList = doctorIDs
	guild.Mafia.CitizenList = citizenIDs
}

// 플레이어 역할 배정
func sendPlayersStartMessage(s *dgo.Session, players []string, numMafia int, numPolice int, numDoctor int) (mafiaIDs []string, policeIDs []string, doctorIDs []string, citizenIDs []string) {
	// embed for Mafia
	embedMafia := dgo.MessageEmbed{
		Title:       "당신은 **마피아**입니다!",
		Description: "밤마다 시민을 처치할 수 있습니다.",
		Color:       0xFFD800,
	}
	// data for Mafia
	dataMafia := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedMafia,
		},
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					&dgo.Button{
						Label:    "Skill",                  // 버튼 텍스트
						Style:    dgo.PrimaryButton,        // 버튼 스타일
						CustomID: "mafia_Game_mafiaButton", // 버튼 클릭 시 처리할 ID
					},
				},
			},
		},
	}

	// embed for Police
	embedPolice := dgo.MessageEmbed{
		Title:       "당신은 **경찰**입니다!",
		Description: "밤마다 한 명의 신원을 확인할 수 있습니다.",
		Color:       0xC87C00,
	}
	// data for Police
	dataPolice := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedPolice,
		},
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					&dgo.Button{
						Label:    "Skill",                  // 버튼 텍스트
						Style:    dgo.PrimaryButton,        // 버튼 스타일
						CustomID: "wink_Game_policeButton", // 버튼 클릭 시 처리할 ID
					},
				},
			},
		},
	}

	// embed for Doctor
	embedDoctor := dgo.MessageEmbed{
		Title:       "당신은 **의사**입니다!",
		Description: "밤마다 마피아로부터 한 명을 보호할 수 있습니다.",
		Color:       0xC87C00,
	}
	// data for Doctor
	dataDoctor := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedDoctor,
		},
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					&dgo.Button{
						Label:    "Skill",                  // 버튼 텍스트
						Style:    dgo.PrimaryButton,        // 버튼 스타일
						CustomID: "wink_Game_doctorButton", // 버튼 클릭 시 처리할 ID
					},
				},
			},
		},
	}

	//embed for Citizen
	embedCitizen := dgo.MessageEmbed{
		Title:       "당신은 **시민**입니다!",
		Description: "마피아를 찾아내세요.",
		Color:       0xC87C00,
	}
	// data for Citizen
	dataCitizen := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedCitizen,
		},
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					&dgo.Button{
						Label:    "Vote",                    // 버튼 텍스트
						Style:    dgo.PrimaryButton,         // 버튼 스타일
						CustomID: "wink_Game_citizenButton", // 버튼 클릭 시 처리할 ID
					},
				},
			},
		},
	}

	shuffled := make([]string, len(players))
	copy(shuffled, players)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	mafiaIDs = shuffled[:numMafia]
	shuffled = shuffled[numMafia:]
	policeIDs = shuffled[:numPolice]
	shuffled = shuffled[numPolice:]
	doctorIDs = shuffled[:numDoctor]
	citizenIDs = shuffled[numDoctor:]

	for _, id := range players {
		switch {
		case general.Contains(mafiaIDs, id):
			err := general.SendComplexDM(s, id, &dataMafia)
			if err != nil {
				log.Printf("Failed to send DM to user %s: %v\n", id, err)
			}
		case general.Contains(policeIDs, id):
			err := general.SendComplexDM(s, id, &dataPolice)
			if err != nil {
				log.Printf("Failed to send DM to user %s: %v\n", id, err)
			}
		case general.Contains(doctorIDs, id):
			err := general.SendComplexDM(s, id, &dataDoctor)
			if err != nil {
				log.Printf("Failed to send DM to user %s: %v\n", id, err)
			}
		default:
			err := general.SendComplexDM(s, id, &dataCitizen)
			if err != nil {
				log.Printf("Failed to send DM to user %s: %v\n", id, err)
			}
		}
	}
	return
}

// func Game_FollowUpMessage(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
// 	msg, err := s.FollowupMessageCreate(i.Interaction, true, &dgo.WebhookParams{
// 		Embeds: []*dgo.MessageEmbed{
// 			{
// 				Title:       "게임은 시작되었습니다",
// 				Description: "마피아 게임이 시작되었습니다! 역할이 개별 DM으로 전송되었습니다.\n",
// 				Color:       0xFFFFFF,
// 			},
// 		},
// 		Components: []dgo.MessageComponent{
// 			dgo.ActionsRow{
// 				Components: []dgo.MessageComponent{
// 					&dgo.Button{
// 						Label:    "게임 재시작",
// 						Style:    dgo.SuccessButton,
// 						CustomID: "wink_restart",
// 					},
// 					&dgo.Button{
// 						Label:    "게임 종료",
// 						Style:    dgo.DangerButton,
// 						CustomID: "wink_end",
// 					},
// 				},
// 			},
// 		},
// 	})

// 	if err != nil {
// 		log.Printf("Failed sending follow-up message [%v]", err)
// 	}

// 	guild.Mafia.MessageIDMap[i.GuildID] = msg.ID
// }

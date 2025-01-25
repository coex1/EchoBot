package mafia

import (

	// system
	"log"
	"math/rand"
	"time"

	// external package
	dgo "github.com/bwmarrin/discordgo"

	// internal package
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"
)

// 플레이어 역할 배정
func sendPlayersStartMessage(s *dgo.Session, guild *data.Guild, players []string, numMafia int, numPolice int, numDoctor int) (mafiaIDs []string, policeIDs []string, doctorIDs []string, citizenIDs []string) {

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
						Label:    "Skill",             // 버튼 텍스트
						Style:    dgo.PrimaryButton,   // 버튼 스타일
						CustomID: "mafia_Skill_Mafia", // 버튼 클릭 시 처리할 ID
					},
					&dgo.Button{
						Label:    "Vote",              // 버튼 텍스트
						Style:    dgo.DangerButton,    // 버튼 스타일
						CustomID: "mafia_Vote_Button", // 버튼 클릭 시 처리할 ID
					},
				},
			},
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					dgo.SelectMenu{
						MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
						CustomID:    "mafia_Alive_listUpdate",
						Placeholder: "생존자 목록",
						MaxValues:   1,
						Options:     guild.Mafia.AllUserInfo,
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
						Label:    "Skill",              // 버튼 텍스트
						Style:    dgo.PrimaryButton,    // 버튼 스타일
						CustomID: "mafia_Skill_Police", // 버튼 클릭 시 처리할 ID
					},
					&dgo.Button{
						Label:    "Vote",              // 버튼 텍스트
						Style:    dgo.DangerButton,    // 버튼 스타일
						CustomID: "mafia_Vote_Button", // 버튼 클릭 시 처리할 ID
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
						Label:    "Skill",               // 버튼 텍스트
						Style:    dgo.PrimaryButton,     // 버튼 스타일
						CustomID: "wmafia_Skill_Doctor", // 버튼 클릭 시 처리할 ID
					},
					&dgo.Button{
						Label:    "Vote",              // 버튼 텍스트
						Style:    dgo.DangerButton,    // 버튼 스타일
						CustomID: "mafia_Vote_Button", // 버튼 클릭 시 처리할 ID
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
						Label:    "Skill",               // 버튼 텍스트
						Style:    dgo.PrimaryButton,     // 버튼 스타일
						CustomID: "mafia_Skill_Citizen", // 버튼 클릭 시 처리할 ID
					},
					&dgo.Button{
						Label:    "Vote",              // 버튼 텍스트
						Style:    dgo.PrimaryButton,   // 버튼 스타일
						CustomID: "mafia_Vote_Button", // 버튼 클릭 시 처리할 ID
					},
				},
			},
		},
	}

	// shuffle algorithm
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

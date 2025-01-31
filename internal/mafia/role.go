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
func sendStartMessage(s *dgo.Session, guild *data.Guild, players []string, numMafia int, numPolice int, numDoctor int) (mafiaIDs []string, policeIDs []string, doctorIDs []string, citizenIDs []string) {

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
						Label:    "준비 완료",
						Style:    dgo.PrimaryButton,
						CustomID: "mafia_Rdy_Button",
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
						Label:    "준비 완료",
						Style:    dgo.PrimaryButton,
						CustomID: "mafia_Rdy_Button",
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
						Label:    "준비 완료",
						Style:    dgo.PrimaryButton,
						CustomID: "mafia_Rdy_Button",
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
						Label:    "준비 완료",
						Style:    dgo.PrimaryButton,
						CustomID: "mafia_Rdy_Button",
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

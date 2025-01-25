package mafia

import (
	"log"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
)

func Game_StartMessage(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	startMessage, err := s.ChannelMessageSendComplex(i.ChannelID, &dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			{
				Title: "게임이 시작되었습니다!",
				Description: "마피아 게임이 시작되었습니다!\n\n" +
					"역할과 진행은 **개별 DM**을 확인해주세요.",
				Color: 0xFFFFFF,
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
		log.Printf("Failed sending follow-up message [%v]", err)
	}
	guild.Mafia.MessageIDMap = make(map[string]string)
	guild.Mafia.MessageIDMap[i.GuildID] = startMessage.ID
}

// ---------- TODO --------------
// func Game_timeUpdate(s *dgo.Session, i *dgo.InteractionCreate, players []string) {
// 	duration := 5 * time.Minute // 5분 타이머
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	endTime := time.Now().Add(duration)

// 	// 남은 시간을 업데이트할 메시지 생성
// 	initialMessage, err := s.ChannelMessageSend(i.ChannelID, "회의 종료까지 5분 남았습니다.")
// 	if err != nil {
// 		log.Printf("Failed to send initial game timer message: %v\n", err)
// 		return
// 	}
// 	// 타이머 루프
// 	go func() {
// 		for range ticker.C {
// 			remaining := time.Until(endTime)
// 			if remaining <= 0 {
// 				// 타이머 종료: 최종 메시지 업데이트
// 				_, err := s.ChannelMessageEdit(initialMessage.ChannelID, initialMessage.ID, "시간이 다 되었습니다! 투표해주세요.")
// 				if err != nil {
// 					log.Printf("Failed to update final game timer message: %v\n", err)
// 				}
// 				break
// 			}

// 			// 남은 시간 계산 및 메시지 업데이트
// 			remainingMinutes := int(remaining.Minutes())
// 			remainingSeconds := int(remaining.Seconds()) % 60
// 			message := fmt.Sprintf("회의 종료까지 %d분 %d초 남았습니다.", remainingMinutes, remainingSeconds)

// 			_, err := s.ChannelMessageEdit(i.ChannelID, initialMessage.ID, message)
// 			if err != nil {
// 				log.Printf("Failed to update game timer message: %v\n", err)
// 				break
// 			}
// 		}
// 	}()
// }

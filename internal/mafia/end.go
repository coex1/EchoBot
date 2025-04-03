package mafia

import (
	"fmt"
	"log"

	dgo "github.com/bwmarrin/discordgo"

	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"
)

func isGameEnd(guild *data.Guild) bool {
	var mafiaNum int
	var citizenNum int

	for _, p := range guild.Mafia.Players {
		if p.IsAlive {
			if p.Role == "Mafia" {
				mafiaNum += 1
			} else {
				citizenNum += 1
			}
		}
	}

	return mafiaNum >= citizenNum
}

func gameEndingMessage(s *dgo.Session, guild *data.Guild) {
	channelID := guild.Mafia.ChannelID
	var mafiaNum int
	var citizenNum int
	var WinnerTeam string
	roleFields := []*dgo.MessageEmbedField{}

	for _, p := range guild.Mafia.Players {
		if p.IsAlive {
			if p.Role == "Mafia" {
				mafiaNum += 1
			} else {
				citizenNum += 1
			}
		}
	}

	if mafiaNum >= citizenNum {
		WinnerTeam = "마피아"
	} else {
		WinnerTeam = "시민"
	}

	for _, p := range guild.Mafia.Players {
		roleFields = append(roleFields, &dgo.MessageEmbedField{
			Name:   p.GlobalName,
			Value:  fmt.Sprintf("(**%s**)", p.Role),
			Inline: true,
		})
	}
	// 투표 내역 임베드
	roleEmbed := &dgo.MessageEmbed{
		Title:       "역할 공개",
		Description: "모든 참여자들의 역할을 확인하세요.",
		Color:       0x3498db,
		Fields:      roleFields,
	}

	for _, p := range guild.Mafia.Players {
		if p.Role == "Mafia" {
			endingEbeds := &dgo.MessageEmbed{
				Title:       fmt.Sprintf("**%s**팀 승리!!", WinnerTeam),
				Description: "승리하였습니다!",
				Color:       0xFFFFFF,
				Fields:      roleFields,
			}
			general.SendComplexDM(s, p.ID, &dgo.MessageSend{
				Embeds: []*dgo.MessageEmbed{endingEbeds, roleEmbed},
			})
		} else {
			endingEbeds := &dgo.MessageEmbed{
				Title:       fmt.Sprintf("**%s**팀 승리!!", WinnerTeam),
				Description: "패배하였습니다.",
				Color:       0xFFFFFF,
				Fields:      roleFields,
			}
			general.SendComplexDM(s, p.ID, &dgo.MessageSend{
				Embeds: []*dgo.MessageEmbed{endingEbeds, roleEmbed},
			})
		}
	}

	// 전체 공지
	_, err := s.ChannelMessageSendComplex(channelID, &dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("**%s**팀 승리!!", WinnerTeam),
				Description: "게임이 종료되었습니다.",
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
					// &dgo.Button{
					// 	Label:    "게임 종료",
					// 	Style:    dgo.DangerButton,
					// 	CustomID: "mafia_end",
					// },
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to send DM to users [%v]", err)
	}
}

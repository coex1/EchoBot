package mafia

import (
	"fmt"
	"log"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
)

func Role_Message(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
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

func Start_Message(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	_, err := s.ChannelMessageSendComplex(i.ChannelID, &dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			{
				Title:       "게임이 시작되었습니다!",
				Description: "좋은 아침입니다!.",
				Color:       0xFFFFFF,
			},
		},
	})
	if err != nil {
		log.Printf("Failed to send DM to users [%v]", err)
	}
}

func Vote_Message(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	var options []dgo.SelectMenuOption
	for _, id := range guild.Mafia.AliveUsersID {
		member, err := s.GuildMember(guild.ID, id)
		if err != nil {
			log.Printf("Failed to get member info for user %s: %v\n", id, err)
			continue
		}
		options = append(options, dgo.SelectMenuOption{
			Label: member.User.GlobalName,
			Value: member.User.ID,
		})
	}

	// embed for Vote
	_, err := s.FollowupMessageCreate(event.Interaction, true, &dgo.WebhookParams{
		// embedVote := dgo.MessageEmbed{
		// 	Title:       "투표",
		// 	Description: "투표할 대상을 선택한 후 '투표하기' 버튼을 눌러주세요!",
		// 	Color:       0xC87C00,
		// }
		//dataVote := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			{
				Title:       "투표",
				Description: "투표할 대상을 선택한 후 '투표하기' 버튼을 눌러주세요!",
				Color:       0xC87C00,
			},
		},
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					&dgo.SelectMenu{
						MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
						CustomID:    "mafia_Vote_listUpdate",
						Placeholder: "한 명을 선택하세요",
						MaxValues:   1,
						Options:     options,
					},
				},
			},
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					&dgo.Button{
						Label:    "투표하기",
						Style:    dgo.PrimaryButton,
						CustomID: "mafia_Vote_Button",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to FollowUpMessage [%v]", err)
	}
	// for _, id := range guild.Mafia.AliveUsersID {
	// 	err := general.SendComplexDM(s, id, &dataVote)
	// 	if err != nil {
	// 		log.Printf("Fauiled to send DM to user %s: %v\n", id, err)
	// 	}
	// }
}

func Vote_listUpdate(i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.VoteMap[i.User.ID] = i.MessageComponentData().Values[0]
}

// on interaction event 'mafia_Vote_Button'
func Vote_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {

	// 선택 완료 메시지
	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseUpdateMessage,
		Data: &dgo.InteractionResponseData{
			Content: fmt.Sprintf("%s에게 투표!", guild.Mafia.VoteMap[i.User.ID]),
		},
	})
	if err != nil {
		log.Printf("Failed to send vote confirmation: %v\n", err)
	}
	log.Printf("User %s finalized vote for %s", i.User.GlobalName, guild.Mafia.VoteMap[i.User.ID])
}

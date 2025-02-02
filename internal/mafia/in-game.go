package mafia

import (
	"fmt"
	"log"
	"strconv"

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

func Day_Message(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	day := strconv.Itoa(guild.Mafia.Day)
	for _, userID := range guild.Mafia.SelectedUsersID {
		channel, err := s.UserChannelCreate(userID)
		if err != nil {
			log.Printf("Failed to create DM channel for user %s: %v\n", userID, err)
			continue
		}
		message := &dgo.MessageSend{
			Embeds: []*dgo.MessageEmbed{
				{
					Title:       day + "일 차 아침입니다.",
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
							Options:     guild.Mafia.AliveUsersID,
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
		}
		_, err = s.ChannelMessageSendComplex(channel.ID, message)
		if err != nil {
			log.Printf("Failed to send confirmation DM to user %s: %v\n", userID, err)
		}
	}
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

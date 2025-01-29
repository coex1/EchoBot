package mafia

import (
	"fmt"
	"log"
	"sync"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"
)

func Start_Message(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	_, err := s.ChannelMessageSendComplex(i.ChannelID, &dgo.MessageSend{
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
		log.Printf("Failed to send DM to users [%v]", err)
	}
	// guild.Mafia.MessageIDMap = make(map[string]string)
	// guild.Mafia.MessageIDMap[i.GuildID] = startMessage.ID
}

func Vote_Message(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild, aliveUsers []string) {
	var options []dgo.SelectMenuOption
	for _, playerID := range aliveUsers {
		member, err := s.GuildMember(guild.ID, playerID)
		if err != nil {
			log.Printf("Failed to get member info for user %s: %v\n", playerID, err)
			continue
		}
		options = append(options, dgo.SelectMenuOption{
			Label: member.User.GlobalName,
			Value: member.User.ID,
		})
	}

	// _, err := s.ChannelMessageSendComplex(event.ChannelID, &dgo.MessageSend{
	// embed for Vote
	embedVote := dgo.MessageEmbed{
		Title:       "투표",
		Description: "투표할 대상을 선택한 후 '투표하기' 버튼을 눌러주세요!",
		Color:       0xC87C00,
	}
	dataVote := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedVote,
		},
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					&dgo.SelectMenu{
						CustomID:    "mafia_Vote_Select",
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
						CustomID: "mafia_Vote_Submit",
					},
				},
			},
		},
	}
	for _, id := range aliveUsers {
		err := general.SendComplexDM(s, id, &dataVote)
		if err != nil {
			log.Printf("Fauiled to send DM to user %s: %v\n", id, err)
		}
	}
}

func Vote_listUpdate(i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.AliveUsersID = i.MessageComponentData().Values
}

// on interaction event 'mafia_Vote_Button'
func Vote_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	var voteMutex = &sync.Mutex{}

	selectedUserID := i.MessageComponentData().Values[0] // 선택된 유저 ID
	member, _ := s.GuildMember(guild.ID, selectedUserID)
	selectedUserGN := member.User.GlobalName // 선택된 유저 Global Name
	voterID := i.User.GlobalName             // 투표한 유저 ID

	// 동기화하여 voteMap에 저장
	voteMutex.Lock()
	guild.Mafia.VoteMap[voterID] = selectedUserGN
	voteMutex.Unlock()

	log.Printf("User %s voted for %s", voterID, selectedUserGN)

	// 선택 완료 메시지
	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseUpdateMessage,
		Data: &dgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ %s에게 투표하셨습니다!", selectedUserGN),
		},
	})
	if err != nil {
		log.Printf("Failed to send vote confirmation: %v\n", err)
	}
}

func Vote_Submit(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	var voteMutex = &sync.Mutex{}

	voterID := i.User.ID

	voteMutex.Lock()
	selectedUser, exists := guild.Mafia.VoteMap[voterID]
	voteMutex.Unlock()

	if !exists {
		err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionResponseUpdateMessage,
			Data: &dgo.InteractionResponseData{
				Content: "투표할 대상을 선택한 후 다시 시도해주세요.",
			},
		})
		if err != nil {
			log.Printf("Failed to send vote warning: %v\n", err)
		}
		return
	}

	// 투표 완료 메시지
	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseUpdateMessage,
		Data: &dgo.InteractionResponseData{
			Content: fmt.Sprintf("🗳️ 최종 투표: %s에게 투표 완료!", selectedUser),
		},
	})
	if err != nil {
		log.Printf("Failed to send vote confirmation: %v\n", err)
	}

	log.Printf("User %s finalized vote for %s", voterID, selectedUser)
}

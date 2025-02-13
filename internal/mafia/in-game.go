package mafia

import (
	"fmt"
	"log"
	"strconv"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"
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

func Day_Message(s *dgo.Session, guild *data.Guild) {
	// 매일 VoteMap, VoteCount 초기화
	guild.Mafia.VoteMap = make(map[string]string)
	guild.Mafia.VoteCount = make(map[string]int)

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
							Options:     guild.Mafia.AliveUsersIDMap,
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

// on interaction event 'mafia_Vote_Button'rintf("%v", guild.Mafia.VoteCount)
func Vote_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	allVoteCount := func(voteMap map[string]int) int {
		total := 0
		for _, cnt := range voteMap {
			total += cnt
		}
		return total
	}
	userID := i.User.ID
	votedPlayer := guild.Mafia.VoteMap[userID]

	// 투표 집계 업데이트
	guild.Mafia.VoteCount[votedPlayer] += 1

	// 모든 플레이어가 투표를 완료했는지 확인
	if allVoteCount(guild.Mafia.VoteCount) == len(guild.Mafia.SelectedUsersID) {
		log.Printf("All players have voted (%d). Sending results...", allVoteCount(guild.Mafia.VoteCount))

		// 투표 결과 집계
		var maxVotes int
		var eliminatedPlayer string
		voteFields := []*dgo.MessageEmbedField{}  // 개별 투표 내역
		countFields := []*dgo.MessageEmbedField{} // 총 투표 집계

		for voter, voted := range guild.Mafia.VoteMap {
			voterName := guild.Mafia.Players[voter].GlobalName

			// 누가 누구에게 투표했는지 추가
			voteFields = append(voteFields, &dgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s ", voterName),
				Value:  fmt.Sprintf("-->> <@%s>", voted),
				Inline: true,
			})
		}

		for player, votes := range guild.Mafia.VoteCount {
			// 총 투표 개수 추가
			countFields = append(countFields, &dgo.MessageEmbedField{
				Name:   fmt.Sprintf("<@%s> ", guild.Mafia.Players[player].GlobalName),
				Value:  fmt.Sprintf("%d 표", votes),
				Inline: true,
			})

			// 가장 많이 득표한 플레이어 찾기
			if votes > maxVotes {
				maxVotes = votes
				eliminatedPlayer = guild.Mafia.Players[player].GlobalName
			}
		}

		// 최다 득표 플레이어 메시지 추가
		finalMessage := fmt.Sprintf("가장 많은 표를 받은 플레이어는 **<@%s>** 입니다!", eliminatedPlayer)

		// 투표 내역 임베드
		voteEmbed := &dgo.MessageEmbed{
			Title:       "투표 내역",
			Description: "누가 누구에게 투표했는지 확인하세요.",
			Color:       0x3498db,
			Fields:      voteFields,
			Footer: &dgo.MessageEmbedFooter{
				Text: "투표가 종료되었습니다.",
			},
		}

		// 최종 투표 결과 임베드
		resultEmbed := &dgo.MessageEmbed{
			Title:       "투표 결과",
			Description: finalMessage,
			Color:       0xe74c3c, // 빨간색
			Fields:      countFields,
			Footer: &dgo.MessageEmbedFooter{
				Text: "투표 결과를 바탕으로 게임을 진행하세요.",
			},
		}

		// 모든 플레이어에게 DM으로 결과 전송
		for _, playerID := range guild.Mafia.AliveUsersID {
			general.SendComplexDM(s, playerID, &dgo.MessageSend{
				Embeds: []*dgo.MessageEmbed{voteEmbed, resultEmbed},
			})
		}

		log.Println("Vote results sent to all players.")

		guild.Mafia.VoteMap = make(map[string]string)

		// 게임 진행 관련 다음 단계 로직 추가 가능 (예: 밤 시간으로 전환)
		guild.Mafia.State = false
	}
}

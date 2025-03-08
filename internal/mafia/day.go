package mafia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"
)

func Day_Message(s *dgo.Session, guild *data.Guild) {
	AliveUsersID := Vote_Reset(guild) // return AliveUserID []dgo.SelectMenuOption
	day := strconv.Itoa(guild.Mafia.Day)

	for _, player := range guild.Mafia.Players {
		_, err := s.UserChannelCreate(player.ID)
		if err != nil {
			log.Printf("Failed to create DM channel for user %s: %v\n", player.ID, err)
			continue
		}
		message := &dgo.MessageSend{
			Embeds: []*dgo.MessageEmbed{
				{
					Title:       day + "일 차 아침입니다.",
					Description: "10 분 동안 토론하며 투표할 대상을 선택해 주세요!",
					Color:       0xC87C00,
				},
			},
			// 투표 버튼
			Components: []dgo.MessageComponent{
				dgo.ActionsRow{
					Components: []dgo.MessageComponent{
						&dgo.SelectMenu{
							MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
							CustomID:    "mafia_Vote_listUpdate",
							Placeholder: "한 명을 선택하세요",
							MaxValues:   1,
							Options:     AliveUsersID,
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
		_, err = s.ChannelMessageSendComplex(player.DMChannelID, message)
		if err != nil {
			log.Printf("Failed to send confirmation DM to user %s: %v\n", player.ID, err)
		}
	}
}

func oneMinute_left(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {

}

func announce_results(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {

}

func Vote_listUpdate(i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.VoteMap[i.User.ID] = i.MessageComponentData().Values[0]
}

// on interaction event 'mafia_Vote_Button'
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
	numAlive := 0
	for _, player := range guild.Mafia.Players {
		if player.IsAlive {
			numAlive += 1
		}
	}

	// 투표 집계 업데이트
	guild.Mafia.VoteCount[votedPlayer] += 1

	// 모든 플레이어가 투표를 완료했을 경우
	if allVoteCount(guild.Mafia.VoteCount) == numAlive {
		log.Printf("All players have voted (%d). Sending results...", allVoteCount(guild.Mafia.VoteCount))

		// 투표 결과 집계
		var maxVotes int
		var SelectedPlayer string
		var SelectedPlayerID string
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

		for id, votes := range guild.Mafia.VoteCount {
			// 총 투표 개수 추가
			countFields = append(countFields, &dgo.MessageEmbedField{
				Name:   fmt.Sprintf("<@%s>\n", guild.Mafia.Players[id].GlobalName),
				Value:  fmt.Sprintf("%d 표", votes),
				Inline: true,
			})

			// 가장 많이 득표한 플레이어 찾기
			if votes > maxVotes {
				maxVotes = votes
				SelectedPlayer = guild.Mafia.Players[id].GlobalName
				SelectedPlayerID = id
			}
		}

		guild.Mafia.Players[SelectedPlayerID].IsAlive = false

		// 최다 득표 플레이어 메시지 추가
		finalMessage := fmt.Sprintf("<@%s>**님이 투표로 처형되었습니다.", SelectedPlayer)

		// 투표 내역 임베드
		voteEmbed := &dgo.MessageEmbed{
			Title:       "투표 내역",
			Description: "누가 누구에게 투표했는지 확인하세요.",
			Color:       0x3498db,
			Fields:      voteFields,
		}

		// 최종 투표 결과 임베드
		resultEmbed := &dgo.MessageEmbed{
			Title:       "투표 결과",
			Description: finalMessage,
			Color:       0xe74c3c, // 빨간색
			Fields:      countFields,
		}

		// 모든 플레이어에게 DM으로 결과 전송
		for _, player := range guild.Mafia.Players {
			general.SendComplexDM(s, player.ID, &dgo.MessageSend{
				Embeds: []*dgo.MessageEmbed{voteEmbed, resultEmbed},
			})
		}
		log.Println("Vote results sent to all players.")

		guild.Mafia.State = false
	} else {
		oneMinute_left(s, i, guild)
		go func() {
			time.Sleep(5 * time.Second) // TODO: change to 15
			announce_results(s, i, guild)
		}()
	}
}

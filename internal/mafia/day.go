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

func Day_Message(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	// 투표 정보 초기화
	players := guild.Mafia.Players
	AliveUsersID := Reset(guild) // return AliveUserID []dgo.SelectMenuOption

	// 일 수 + 1
	guild.Mafia.Day += 1
	day := strconv.Itoa(guild.Mafia.Day)

	for _, p := range players {
		if p.IsAlive {
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
							&dgo.Button{
								Label:    "기권하기",
								Style:    dgo.SecondaryButton,
								CustomID: "mafia_Skip_Button",
							},
						},
					},
				},
			}
			general.SendComplexDM(s, p.ID, message)
		}
	}

	// 10분 후 자동 기권 처리
	guild.Mafia.TimerActive = true

	go func() {
		time.Sleep(10 * time.Minute)
		if guild.Mafia.TimerActive {
			autoSkipUnvotedPlayers(guild)

			announceVoteResult(s, i, guild)
		}
	}()
}

func autoSkipUnvotedPlayers(guild *data.Guild) {
	for playerID := range guild.Mafia.Players {
		if _, voted := guild.Mafia.VoteMap[playerID]; voted {
			continue
		}
		if general.Contains(guild.Mafia.VoteSkip, playerID) {
			continue
		}

		guild.Mafia.VoteSkip = append(guild.Mafia.VoteSkip, playerID)
		log.Printf("User %s did not vote in time, automatically skipped", playerID)
	}
}

// func oneMinute_left(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
// 	// TODO
// }

func Vote_listUpdate(i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.TempVoteMap[i.User.ID] = i.MessageComponentData().Values[0]
}

// on interaction event 'mafia_Vote_Button'
func Vote_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	for id := range guild.Mafia.VoteMap {
		if id == i.User.ID {
			return
		}
	}
	// 투표 집계 업데이트
	if vote, exists := guild.Mafia.TempVoteMap[i.User.ID]; exists {
		guild.Mafia.VoteMap[i.User.ID] = vote
		guild.Mafia.VoteCount[vote] += 1
	}
	log.Printf("User %s Vote %s", i.User.ID, guild.Mafia.VoteMap[i.User.ID])

	announceVoteResult(s, i, guild)
}

func Skip_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	for _, skipped := range guild.Mafia.VoteSkip {
		if skipped == i.User.ID {
			return
		}
	}
	guild.Mafia.VoteSkip = append(guild.Mafia.VoteSkip, i.User.ID)
	delete(guild.Mafia.TempVoteMap, i.User.ID)
	delete(guild.Mafia.VoteMap, i.User.ID)
	log.Printf("User %s skipped voting", guild.Mafia.Players[i.User.ID].GlobalName)

	announceVoteResult(s, i, guild)
}

func announceVoteResult(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	numAlive := 0
	for _, player := range guild.Mafia.Players {
		if player.IsAlive {
			numAlive += 1
		}
	}
	// 모든 플레이어가 투표(or 기권)를 완료했을 경우
	guild.Mafia.TimerActive = false

	if len(guild.Mafia.VoteMap)+len(guild.Mafia.VoteSkip) == numAlive {
		log.Printf("All players have voted (%d). Sending results...", numAlive)
		majority := numAlive / 2

		// === 투표 결과 정리 ===
		var maxVotes int
		var selectedPlayerID string
		voteFields := []*dgo.MessageEmbedField{}  // 개별 투표 내역
		countFields := []*dgo.MessageEmbedField{} // 투표 결과
		voteCounts := make(map[int][]string)      // 동점 체크용 (표 수 : 플레이어 목록)

		// 누가 누구에게 투표했는지 추가
		for voter, voted := range guild.Mafia.VoteMap {
			voteFields = append(voteFields, &dgo.MessageEmbedField{
				Name:   guild.Mafia.Players[voter].GlobalName,
				Value:  fmt.Sprintf("-->> **%s**", guild.Mafia.Players[voted].GlobalName),
				Inline: true,
			})
		}
		// 누가 스킵했는지 추가
		for _, id := range guild.Mafia.VoteSkip {
			voteFields = append(voteFields, &dgo.MessageEmbedField{
				Name:   guild.Mafia.Players[id].GlobalName,
				Value:  "(**기권**)",
				Inline: true,
			})
		}
		for id, votes := range guild.Mafia.VoteCount {
			// 총 투표 개수 추가
			countFields = append(countFields, &dgo.MessageEmbedField{
				Name:   guild.Mafia.Players[id].GlobalName,
				Value:  fmt.Sprintf("**%d 표**", votes),
				Inline: true,
			})
			// 가장 많이 득표한 플레이어 찾기
			if votes > maxVotes {
				maxVotes = votes
				selectedPlayerID = id
			}
			voteCounts[votes] = append(voteCounts[votes], id)
		}

		// 투표 내역 임베드
		voteEmbed := &dgo.MessageEmbed{
			Title:       "투표 내역",
			Description: "누가 누구에게 투표했는지 확인하세요.",
			Color:       0x3498db,
			Fields:      voteFields,
		}

		// 케이스별 처리
		if len(guild.Mafia.VoteMap) == 0 {
			// 전원 기권
			resultEmbed := &dgo.MessageEmbed{
				Title:       "투표 결과",
				Description: "모든 플레이어가 기권했습니다. 아무도 처형되지 않았습니다.",
				Color:       0xe74c3c,
			}
			for _, player := range guild.Mafia.Players {
				general.SendComplexDM(s, player.ID, &dgo.MessageSend{
					Embeds: []*dgo.MessageEmbed{voteEmbed, resultEmbed},
				})
			}
			log.Println("All players skipped voting.")
		} else if len(voteCounts[maxVotes]) > 1 {
			// 과반수를 넘지 않거나 동점
			resultEmbed := &dgo.MessageEmbed{
				Title:       "투표 결과",
				Description: "아무도 처형되지 않았습니다.",
				Color:       0xe74c3c,
				Fields:      countFields,
			}
			// 모든 플레이어에게 DM으로 결과 전송
			for _, player := range guild.Mafia.Players {
				general.SendComplexDM(s, player.ID, &dgo.MessageSend{
					Embeds: []*dgo.MessageEmbed{voteEmbed, resultEmbed},
				})
			}
		} else if maxVotes >= majority {
			selectedPlayer := guild.Mafia.Players[selectedPlayerID].GlobalName
			guild.Mafia.Players[selectedPlayerID].IsAlive = false
			// 과반수를 넘음
			resultEmbed := &dgo.MessageEmbed{
				Title:       "투표 결과",
				Description: fmt.Sprintf("**%s** 님이 과반수로 처형되었습니다.", selectedPlayer),
				Color:       0xe74c3c, // 빨간색
				Fields:      countFields,
			}

			// 모든 플레이어에게 DM으로 결과 전송
			for _, player := range guild.Mafia.Players {
				general.SendComplexDM(s, player.ID, &dgo.MessageSend{
					Embeds: []*dgo.MessageEmbed{voteEmbed, resultEmbed},
				})
			}
		}

		log.Println("Vote results sent to all players.")

		time.Sleep(5 * time.Second)

		if isGameOver(guild) {
			gameEndingMessage(s, i, guild)
		} else {
			Night_Message(s, guild)
		}
	}
}

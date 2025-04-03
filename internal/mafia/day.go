package mafia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"
)

func Day_Message(s *dgo.Session, guild *data.Guild) {
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

	StartVoteTimer(s, guild)
}

func StartVoteTimer(s *dgo.Session, guild *data.Guild) {
	// 이전 타이머 취소
	if guild.Mafia.CancelFunc != nil {
		guild.Mafia.CancelFunc()
	}

	// 10 분 타이머
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	guild.Mafia.CancelFunc = cancel
	guild.Mafia.TimerActive = true

	// 종료 1분 전 타이머
	go func() {
		select {
		case <-time.After(9 * time.Minute):
			if ctx.Err() == nil && guild.Mafia.TimerActive {
				oneMinute_left(s, guild)
			}
		case <-ctx.Done():
		}
	}()

	// 10분 타이머 종료 처리
	go func() {
		<-ctx.Done()

		if !guild.Mafia.TimerActive {
			log.Println("Ignore goroutine")
			return
		}

		if ctx.Err() == context.DeadlineExceeded {
			log.Println("AutoSkip Voting!")
			autoSkipUnvotedPlayers(guild)
			announceVoteResult(s, guild)
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

func oneMinute_left(s *dgo.Session, guild *data.Guild) {
	for _, p := range guild.Mafia.Players {
		if p.IsAlive {
			general.SendDM(s, p.ID, "투표 종료까지 **1분** 남았습니다!")
		}
	}
}

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

	announceVoteResult(s, guild)
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

	announceVoteResult(s, guild)
}

func announceVoteResult(s *dgo.Session, guild *data.Guild) {
	numAlive := 0
	for _, player := range guild.Mafia.Players {
		if player.IsAlive {
			numAlive += 1
		}
	}

	if len(guild.Mafia.VoteMap)+len(guild.Mafia.VoteSkip) == numAlive {
		log.Printf("All players have voted (%d). Sending results...", numAlive)

		// 고루틴 종료
		guild.Mafia.TimerActive = false
		if guild.Mafia.CancelFunc != nil {
			guild.Mafia.CancelFunc()
		}

		// === 투표 결과 정리 ===
		var maxVotes int
		var selectedPlayerID string
		voteFields := []*dgo.MessageEmbedField{}  // 개별 투표 내역
		countFields := []*dgo.MessageEmbedField{} // 투표 결과
		voteCounts := make(map[int][]string)      // 동점 체크용 (표 수 : 플레이어 목록)

		// 투표 필드
		for voter, voted := range guild.Mafia.VoteMap {
			voteFields = append(voteFields, &dgo.MessageEmbedField{
				Name:   guild.Mafia.Players[voter].GlobalName,
				Value:  fmt.Sprintf("-->> **%s**", guild.Mafia.Players[voted].GlobalName),
				Inline: true,
			})
		}
		// 스킵 필드
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
			voteCounts[votes] = append(voteCounts[votes], id)
			if votes > maxVotes {
				maxVotes = votes
				selectedPlayerID = id
			}
		}

		// 투표 내역 임베드
		voteEmbed := &dgo.MessageEmbed{
			Title:       "투표 내역",
			Description: "누가 누구에게 투표했는지 확인하세요.",
			Color:       0x3498db,
			Fields:      voteFields,
		}

		var resultEmbed *dgo.MessageEmbed

		// 케이스별 처리
		if len(guild.Mafia.VoteMap) == 0 {
			// 전원 기권
			resultEmbed = &dgo.MessageEmbed{
				Title:       "투표 결과",
				Description: "모든 플레이어가 기권했습니다. 아무도 처형되지 않았습니다.",
				Color:       0xe74c3c,
			}
			log.Println("All players skipped voting.")
		} else if len(voteCounts[maxVotes]) > 1 {
			// 과반수를 넘지 않거나 동점
			resultEmbed = &dgo.MessageEmbed{
				Title:       "투표 결과",
				Description: "동점으로 아무도 처형되지 않았습니다.",
				Color:       0xe74c3c,
				Fields:      countFields,
			}
		} else if maxVotes > numAlive/2 {
			selectedPlayer := guild.Mafia.Players[selectedPlayerID].GlobalName
			guild.Mafia.Players[selectedPlayerID].IsAlive = false
			// 과반수를 넘음
			resultEmbed = &dgo.MessageEmbed{
				Title:       "투표 결과",
				Description: fmt.Sprintf("**%s** 님이 과반수로 처형되었습니다.", selectedPlayer),
				Color:       0xe74c3c, // 빨간색
				Fields:      countFields,
			}
		} else {
			// 과반수를 넘지 않거나 동점
			resultEmbed = &dgo.MessageEmbed{
				Title:       "투표 결과",
				Description: "과반수를 넘지 않아 아무도 처형되지 않았습니다.",
				Color:       0xe74c3c,
				Fields:      countFields,
			}
		}

		for _, player := range guild.Mafia.Players {
			general.SendComplexDM(s, player.ID, &dgo.MessageSend{
				Embeds: []*dgo.MessageEmbed{voteEmbed, resultEmbed},
			})
		}
		log.Println("All players skipped voting.")

		log.Println("Vote results sent to all players.")

		time.Sleep(5 * time.Second)

		if isGameEnd(guild) {
			gameEndingMessage(s, guild)
		} else {
			Night_Message(s, guild)
		}
	}
}

package wink

import (
  // system packages
	"fmt"
	"log"
	"strconv"
  "strings"

  // internal packages
  "github.com/coex1/EchoBot/internal/general"
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// on interaction event 'wink_Start_Button'
func Start_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  // check if player count is valid
  players := guild.Wink.SelectedUsers[i.GuildID]
  if len(players) < MIN_PLAYER_CNT {
    log.Printf("Invalid player count, ending game")
    Start_Failed(s, i, guild, "인원수가 부족")
    return
  }
	guild.Wink.TotalParticipants = len(players)

  // select king
  kingID := selectKing(players)

  // send role notice via private DM
  sendPlayersStartMessage(s, i, guild, players, kingID)

  // send FollowUp message
  Game_FollowUpMessage(s, i, guild)
}















func FollowUpHandler(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	userID := event.Member.User.ID
	userGlobalName := event.Member.User.GlobalName
	action := ""

	// 버튼 클릭에 따라 상태 업데이트
	switch event.MessageComponentData().CustomID {
	case "wink_check":
		if !guild.Wink.CheckedUsers[userID] {
			guild.Wink.CheckedUsers[userID] = true
		}
		action = "확인"
	case "wink_cancel":
		if guild.Wink.CheckedUsers[userID] {
			guild.Wink.CheckedUsers[userID] = false
		}
		action = "취소"
	}

	// 현재 체크된 수 계산
	checkedCount := general.CountCheckedUsers(guild.Wink.CheckedUsers)

	// 기존 메시지 업데이트 (Followup 메시지 수정)
	messageID := guild.Wink.MessageIDMap[event.GuildID]
	if messageID == "" {
		log.Println("No message found to update")
		return
	}

	// 체크된 유저 및 체크되지 않은 유저 목록 생성
	var List, uncheckedUsersList string
	for _, id := range guild.Wink.SelectedUsers[event.GuildID] {
		// 유저 정보를 가져오기
		member, err := s.GuildMember(event.GuildID, id)
		if err != nil {
			log.Println("Error fetching member:", err)
			continue
		}
		userName := member.User.GlobalName

		if guild.Wink.CheckedUsers[id] {
			List += fmt.Sprintf("%s\n", userName)
		} else {
			uncheckedUsersList += fmt.Sprintf("%s\n", userName)
		}
	}

	// 남은 사람이 한 명일 경우 처리
	var embed *dgo.MessageEmbed
	if checkedCount == guild.Wink.TotalParticipants-1 {
		lastUserName := uncheckedUsersList
		embed = &dgo.MessageEmbed{
			Title: "마지막 남은 사람!\n",
			Description: fmt.Sprintf(
				"%s님, 당신이 마지막 사람입니다.\n\n왕일 것 같은 사람을 지목해주세요!", strings.ReplaceAll(lastUserName, "\n", ""),
			),
			Color: 0xff0000, // 다른 색으로 표시
		}
	} else {
		embed = &dgo.MessageEmbed{
			Title: "게임 진행 중...\n",
			Description: fmt.Sprintf(
				"윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** %d / %d\n\n**확인한 유저 :**\n%s\n**확인하지 못한 유저 :**\n%s",
				checkedCount, guild.Wink.TotalParticipants, List, uncheckedUsersList,
			),
			Color: 0x00ff00,
		}
	}

	// 상호작용 응답 지연 후 아래 메시지 수정 진행
	err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	// 메시지 수정
	content := fmt.Sprintf("'%s'이(가) %s했습니다.\n", userGlobalName, action)
	_, err = s.ChannelMessageEditComplex(&dgo.MessageEdit{
		Channel:    event.ChannelID,
		ID:         messageID,
		Embeds:     &[]*dgo.MessageEmbed{embed},
		Content:    &content,
		Components: &event.Message.Components, // 기존 버튼 컴포넌트 유지
	})
	if err != nil {
		log.Println("Error updating message:", err)
		return
	}
}

// 버튼 포함 임베드 메시지 생성
func CreateFollowUpMessage(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Wink.TotalParticipants = len(guild.Wink.SelectedUsers[i.GuildID])

	embed := &dgo.MessageEmbed{
		Title:       "게임 시작!",
		Description: "윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** 0 / " + strconv.Itoa(guild.Wink.TotalParticipants),
		Color:       0x00ff00,
	}

	components := []dgo.MessageComponent{
		dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				&dgo.Button{
					Label:    "V",
					Style:    dgo.SuccessButton,
					CustomID: "wink_check",
				},
				&dgo.Button{
					Label:    "X",
					Style:    dgo.DangerButton,
					CustomID: "wink_cancel",
				},
			},
		},
	}

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &dgo.WebhookParams{
		Embeds:     []*dgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		log.Println("Error sending follow-up message:", err)
	}

	guild.Wink.MessageIDMap[i.GuildID] = msg.ID
}

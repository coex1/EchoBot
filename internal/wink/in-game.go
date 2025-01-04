package wink

import (
  // system packages
	"log"
  "fmt"
  "strings"
  "time"

  // internal packages
  "github.com/coex1/EchoBot/internal/general"
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// on interaction event 'wink_Game_listUpdate'
func Game_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  guild.Wink.UserSelection[event.User.GlobalName] = event.MessageComponentData().Values[0]
}

func Game_submitButton(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	if guild.Wink.ConfirmedUsers[i.User.GlobalName] == false {
    target := guild.Wink.UserSelection[i.User.GlobalName]

    guild.Wink.ConfirmedCount++
    guild.Wink.UserSelectionFinal[i.User.GlobalName] = target
    guild.Wink.ConfirmedUsers[i.User.GlobalName] = true
    log.Printf("-----> true cnt = %d", guild.Wink.ConfirmedCount)

    // TODO: upgrade log to be a better log
    log.Printf("[" + i.User.GlobalName +"] selected user [" + target + "]")

    // ignore index
    general.SendDM(s, i.User.ID, "지목하신 상대는 [" + target + "] 입니다!\n(원하시면 언제든지 수정하실 수 있으십니다!)")
  }

  go func() {
    checkEndCondition(s, guild)
  }()
}

func Game_submitFakeButton(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	if guild.Wink.ConfirmedUsers[i.User.GlobalName] == false {
    guild.Wink.ConfirmedCount++
    guild.Wink.ConfirmedUsers[i.User.GlobalName] = true
    log.Printf("-----> true cnt = %d", guild.Wink.ConfirmedCount)
  }

  // TODO: upgrade log to be a better log
  log.Printf("the king has inputted a fake signal")

  // ignore index
  general.SendDM(s, i.User.ID, "윙크 받으셨다고 처리되었습니다!")

  go func() {
    checkEndCondition(s, guild)
  }()
}

func checkEndCondition(s *dgo.Session, guild *data.Guild) {
  final_person_global_name := ""

  if guild.Wink.TotalParticipants-1 == guild.Wink.ConfirmedCount {
    log.Println("map? : ", guild.Wink.ConfirmedUsers)
    for k, i :=	range guild.Wink.ConfirmedUsers {
      log.Printf("---------> iter test = k=%s i=%d", k, i)
      if i == false {
        final_person_global_name = k
      }
    }
    
    log.Printf("Ending game!!! final_person_global_name= [%s]", final_person_global_name)
    send_noti_of_final_person(s, guild, final_person_global_name)
  }
}

func send_noti_of_final_person(s *dgo.Session, guild *data.Guild, f string) {
  if f == guild.Wink.KingName {
    announce_results()
  } else {
    announce_last_person(s, guild, f)
    go func() {
      // start game end timer (default: 15 sec?)
      time.Sleep(15)
      announce_results()
    }()
  }
}

// announce to everyone the last person
func announce_last_person (s *dgo.Session, guild *data.Guild, f string) {
  players := guild.Wink.SelectedUsersID

  embed := dgo.MessageEmbed{
    Title:        "마지막 사람이.....",
    Description:  "\""+f+"\" 입니다!!!\n" +
                  "15초 뒤에 게임이 종료되니 얼른 투표해!",
    Color:        0xFC2803,
  }
  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      &embed,
    },
  }

  // ignore array index
  for _, i := range players {
    general.SendComplexDM(s, i, &data)
  }
}

// announce to everyone game results
func announce_results() {
  
}

// ?
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
	for _, id := range guild.Wink.SelectedUsersID {
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

  // 메시지 수정
  content := fmt.Sprintf("'%s'이(가) %s했습니다.\n", userGlobalName, action)
  _, err := s.ChannelMessageEditComplex(&dgo.MessageEdit{
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

package wink

import (
  // system packages
	"log"

  // internal packages
  "github.com/coex1/EchoBot/internal/general"
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// on interaction event 'wink_Start_listUpdate'
// update selected user list
func Start_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  guild.Wink.SelectedUsersID = event.MessageComponentData().Values
}

// on interaction event 'wink_Start_Button'
func Start_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  players := guild.Wink.SelectedUsersID

  // TODO: run after game state check
  reset_fully(guild)

  // check if player count is valid
  if len(players) < MIN_PLAYER_CNT {
    log.Printf("Invalid player count, ending game")
    Start_Failed(s, i, guild, "인원수가 부족")
    return
  }
	guild.Wink.TotalParticipants = len(players)

  for _, u := range guild.Wink.AllUserInfo {
    isPart := false

    for _, a := range guild.Wink.SelectedUsersID {
      if u.Value == a {
        isPart = true
        break
      }
    }

    if isPart {
      log.Printf("comparing values [%s] [%s]", u.Label, u.Value)
      guild.Wink.SelectedUsersInfo = append(guild.Wink.SelectedUsersInfo, dgo.SelectMenuOption{
        Label: u.Label,
        Value: u.Label,
      })
    }
  }

  // select king
  kingID := selectKing(players)

  // send role notice via private DM
  sendPlayersStartMessage(s, guild, players, kingID)

  // send FollowUp message
  Game_FollowUpMessage(s, i, guild)
}

// on start fail
func Start_Failed(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild, cause string) {
  _, err := s.FollowupMessageCreate(i.Interaction, true, &dgo.WebhookParams{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "게임 시작 실패!",
        Description:  "'" + cause + "' 이유로 게임을 시작 못했습니다 ㅠㅠ",
        Color:        0xFF0000,
      },
    },
  })
  if err != nil {
    log.Printf("Failed sending follow-up message [%v]", err)
	}
}

// 사용자 목록에서 왕 선택
func selectKing(players []string) (kingID string){
	kingID = players[general.Random(0, len(players)-1)]
  log.Printf("Selected king! [%s]", kingID)
  return
}

// 역할 공지 및 선택 메뉴!
// send select menu and confirm button to all users
func sendPlayersStartMessage(s *dgo.Session, guild *data.Guild, players []string, kingID string) {
  var minVal int = 1

  king_embed := dgo.MessageEmbed{
    Title:        "당신은 왕입니다!",
    Description:  "시민 한 사람을 제외한 나머지 사람들에게 윙크를 주세요!\n" +
                  "다른 시민들에게 들키지 않게, 당신도 시민들이 윙크 받았을 때 클릭하는 '제출' 버튼이 있습니다.\n" +
                  "언제든지 윙크 받은 척 하시면서 '제출' 버튼을 클릭 해 주세요.\n" +
                  "(만약 마지막으로 '제출' 버튼을 클릭 하시면 당신의 패배입니다 -.-)\n",
    Color:        0XFFD800,
  }

  villager_embed := dgo.MessageEmbed{
    Title:        "당신은 시민입니다!",
    Description:  "왕으로부터 윙크를 받으세요! (혹은 윙크 하는 것을 발견하세요 ^.-)\n" +
                  "윙크 받으셨으면, 누가 왕인지 목록에서 선택하신 후 '제출' 버튼을 클릭 해 주세요!\n\n" +
                  "참고: 윙크 받으셨으면, 주변 사람들의 눈을 계속 마주쳐 주세요!\n" +
                  "폰을 계속 보고 있으면 뻔히 왕이 아닌것을 눈치치게 될테니....^^;\n",
    Color:        0XC87C00,
  }

  // data for king
  data_k := dgo.MessageSend{
    Components: []dgo.MessageComponent{
      dgo.ActionsRow{
        Components: []dgo.MessageComponent{
          &dgo.Button{
            Label:    "제출",                   // 버튼 텍스트
            Style:    dgo.PrimaryButton,        // 버튼 스타일
            CustomID: "wink_Game_submitFakeButton", // 버튼 클릭 시 처리할 ID
          },
        },
      },
    },
  }

  // data for normies
  data_n := dgo.MessageSend{
    Components: []dgo.MessageComponent{
      dgo.ActionsRow{
        Components: []dgo.MessageComponent{
          dgo.SelectMenu{
            CustomID:     "wink_Game_listUpdate",
            Placeholder:  "사용자 목록",
            MinValues:    &minVal,
            MaxValues:    1,
            Options:      guild.Wink.SelectedUsersInfo,
          },
        },
      },
      dgo.ActionsRow{
        Components: []dgo.MessageComponent{
          &dgo.Button{
            Label:    "제출",                   // 버튼 텍스트
            Style:    dgo.PrimaryButton,        // 버튼 스타일
            CustomID: "wink_Game_submitButton", // 버튼 클릭 시 처리할 ID
          },
        },
      },
    },
  }

  // ignore index
  for _, i := range players {
    if i == kingID {
      // TODO: merge below 2 lines to 1
      data_k.Embeds = []*dgo.MessageEmbed{ &king_embed }
      general.SendComplexDM(s, i, &data_k)
    } else {
      // TODO: merge below 2 lines to 1
      data_n.Embeds = []*dgo.MessageEmbed{ &villager_embed }
      general.SendComplexDM(s, i, &data_n)
    }
  }
}
// send message to guild chat
func Game_FollowUpMessage(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  msg, err := s.FollowupMessageCreate(i.Interaction, true, &dgo.WebhookParams{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "게임은 시작되었습니다",
        Description:  "개인 문자로 역활과 설명을 보냈습니다!\n"+
                      "내용을 읽으시고 게임을 진행하시면 됩니다!",
        Color:        0xFFFFFF,
      },
    },
    Components: []dgo.MessageComponent{
      dgo.ActionsRow{
        Components: []dgo.MessageComponent{
          &dgo.Button{
            Label:    "게임 재시작",
            Style:    dgo.SuccessButton,
            CustomID: "wink_restart",
          },
          &dgo.Button{
            Label:    "게임 종료",
            Style:    dgo.DangerButton,
            CustomID: "wink_end",
          },
        },
      },
    },
  })
  
  if err != nil {
    log.Printf("Failed sending follow-up message [%v]", err)
	}

  guild.Wink.MessageIDMap[i.GuildID] = msg.ID
}


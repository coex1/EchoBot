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

func Start_buttonPressed(s *dgo.Session, i *dgo.InteractionCreate, g *data.Guild) {
  var players []string = g.Wink.SelectedUsersID
  var count int = len(players)
	var list []dgo.SelectMenuOption

  // check if player count is valid
  if len(players) < MIN_PLAYER_CNT {
    log.Printf("Invalid player count, ending game")
    start_sendFailedResponse(s, i, "인원수가 부족")
    return
  }

	g.Wink.TotalParticipants = count

  // create list, for menu, for users to select who they think is the king
  for _, v := range players {
    log.Printf("Player [%s] is included in the game\n", g.Wink.NameList[v])
    
    // initializing arrays
    g.Wink.ConfirmedUsers[v] = false
    g.Wink.UserSelectionFinal[v] = ""

    list = append(list, dgo.SelectMenuOption{
      Label: g.Wink.NameList[v],
      Value: v,
    })
  }

  startGame(s, i, g, list)
}

// on start fail
func start_sendFailedResponse(s *dgo.Session, i *dgo.InteractionCreate, cause string) {
  var embed *dgo.WebhookParams = &dgo.WebhookParams{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "게임 시작 실패!",
        Description:  "'" + cause + "' 이유로 게임을 시작 못했습니다 ㅠㅠ",
        Color:        0xFF0000,
      },
    },
  }

  _, err := s.FollowupMessageCreate(i.Interaction, true, embed)
  if err != nil {
    log.Printf("Failed sending follow-up message [%v]", err)
	}
}

func startGame(s *dgo.Session, i *dgo.InteractionCreate, g *data.Guild, list []dgo.SelectMenuOption) {
  var players []string = g.Wink.SelectedUsersID

  // select king
  g.Wink.KingID = players[general.Random(0, len(players)-1)]
  log.Printf("User [%s] has been selected as this round's king!\n", g.Wink.NameList[g.Wink.KingID])

  // send role notice via private DMs
  sendPlayersRoleNotices(s, g, list)

  // send channel control menu message
  sendChannelControlMenuResponse(s, i)
}

// notify all users of their roles
// and send select menu and confirm button to all users
func sendPlayersRoleNotices(s *dgo.Session, g *data.Guild, list []dgo.SelectMenuOption) {
  var min int = 1

  // data for king
  data_King := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{
      {
        Title:        "당신은 왕입니다!",
        Description:  "시민 한 사람을 제외한 나머지 사람들에게 윙크를 주세요!\n" +
        "다른 시민들에게 들키지 않게, 당신도 시민들이 윙크 받았을 때 클릭하는 '제출' 버튼이 있습니다.\n" +
        "언제든지 윙크 받은 척 하시면서 '제출' 버튼을 클릭 해 주세요.\n" +
        "(만약 마지막으로 '제출' 버튼을 클릭 하시면 당신의 패배입니다 -.-)\n",
        Color:        0xFFD800,
      },
    },
    Components: []dgo.MessageComponent{
      dgo.ActionsRow{
        Components: []dgo.MessageComponent{
          &dgo.Button{
            Label:    "제출",                    // 버튼 텍스트
            Style:    dgo.PrimaryButton,         // 버튼 스타일
            CustomID: "wink_submit_king_button", // 버튼 클릭 시 처리할 ID
          },
        },
      },
    },
  }

  // data for villagers
  data_Norm := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "당신은 시민입니다!",
        Description:  "왕으로부터 윙크를 받으세요! (혹은 윙크 하는 것을 발견하세요 ^.-)\n" +
        "윙크 받으셨으면, 누가 왕인지 목록에서 선택하신 후 '제출' 버튼을 클릭 해 주세요!\n\n" +
        "참고: 윙크 받으셨으면, 주변 사람들의 눈을 계속 마주쳐 주세요!\n" +
        "폰을 계속 보고 있으면 뻔히 왕이 아닌것을 눈치치게 될테니....^^;\n",
        Color:        0xC87C00,
      },
    },
    Components: []dgo.MessageComponent{
      dgo.ActionsRow{
        Components: []dgo.MessageComponent{
          &dgo.SelectMenu{
            MenuType:     dgo.SelectMenuType(dgo.SelectMenuComponent),
            CustomID:     "wink_norm_list",
            Placeholder:  "사용자 목록",
            MinValues:    &min,
            MaxValues:    1,
            Options:      list,
          },
        },
      },
      dgo.ActionsRow{
        Components: []dgo.MessageComponent{
          &dgo.Button{
            Label:    "제출",                    // 버튼 텍스트
            Style:    dgo.PrimaryButton,         // 버튼 스타일
            CustomID: "wink_submit_norm_button", // 버튼 클릭 시 처리할 ID
          },
        },
      },
    },
  }

  // ignore array index
  for _, p := range g.Wink.SelectedUsersID {
    if p == g.Wink.KingID {
      general.SendComplexDM(s, p, &data_King)
    } else {
      general.SendComplexDM(s, p, &data_Norm)
    }
  }
}

// send channel control menu message
func sendChannelControlMenuResponse(s *dgo.Session, i *dgo.InteractionCreate) {
  var embed *dgo.WebhookParams = &dgo.WebhookParams{
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
            Style:    dgo.PrimaryButton,
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
  }

  _, err := s.FollowupMessageCreate(i.Interaction, true, embed)
  if err != nil {
    log.Printf("Failed sending a FollowUp message [%v]", err)
	}
}

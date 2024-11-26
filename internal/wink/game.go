package wink

import (
  // system packages
	"log"
	"strconv"

  // internal packages
  "github.com/coex1/EchoBot/internal/general"
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

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
                  "다른 시민들에게 들키지 않게, 당신도 시민들이 윙크 받았을 때 클릭하는 버튼들이 있습니다.\n" +
                  "(초록 버튼은 받았다, 빨간 버튼은 취소입니다)\n" +
                  "언제든지 윙크 받은 척 하시면서 초록 버튼을 클릭 해 주세요.\n" +
                  "(만약 마지막으로 초록 버튼 클릭 하시면 패배입니다 -.-)\n",
    Color:        0XFFD800,
  }

  villager_embed := dgo.MessageEmbed{
    Title:        "당신은 시민입니다!",
    Description:  "왕으로부터 윙크를 받으세요! (혹은 윙크 하는 것을 발견하세요)\n" +
                  "윙크 받으셨으면 초록 버튼을 클릭 해 주세요!\n" +
                  "(초록 버튼은 받았다, 빨간 버튼은 취소입니다)\n" +
                  "참고: 윙크 받으셨으면 가능한 사람들 눈을 계속 마주쳐 주세요!\n" +
                  "(폰을 계속 보고 있으면 뻔히 왕이 아닌것을 알게 되니....^^;)\n",
    Color:        0XC87C00,
  }

  data := dgo.MessageSend{
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
      data.Embeds = []*dgo.MessageEmbed{ &king_embed }
    } else {
      data.Embeds = []*dgo.MessageEmbed{ &villager_embed }
    }

    general.SendComplexDM(s, i, &data)
  }
}

func Game_FollowUpMessage(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  msg, err := s.FollowupMessageCreate(i.Interaction, true, &dgo.WebhookParams{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "게임 시작!",
        Description:  "윙크를 받으셨으면 V 버튼을 클릭 해 주세요!"+
                      "\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!"+
                      "\n\n**현재 윙크 받은 사람 수 :** 0 / " + strconv.Itoa(guild.Wink.TotalParticipants),
        Color:       0x00ff00,
      },
    },
    Components: []dgo.MessageComponent{
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
    },
  })
  if err != nil {
    log.Printf("Failed sending follow-up message [%v]", err)
	}
  guild.Wink.MessageIDMap[i.GuildID] = msg.ID
}

func checkEndCondition(s *dgo.Session, guild *data.Guild) {
  c := 0
  for _, i :=	range guild.Wink.ConfirmedUsers {
    if i {
      c += 1
    }
  }

  log.Printf("=-----> true cnt = %d", c)
  if guild.Wink.TotalParticipants-1 == c {
    log.Printf("Ending game!!!")
  }

}

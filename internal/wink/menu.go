package wink

import (
  // system packages
	"log"
	"strconv"

  // internal packages
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// Start
var start_selectMenu dgo.SelectMenu = dgo.SelectMenu{
  CustomID:    "wink_Start_listUpdate",
  Placeholder: "사용자를 선택해 주세요!",
}

// on interaction event 'wink_Start_listUpdate'
func Start_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  // update selected user list
  guild.Wink.SelectedUsers[event.GuildID] = event.MessageComponentData().Values

  err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
    // 상호작용 업데이트가 약간의 지연 이후 진행
    Type: dgo.InteractionResponseDeferredMessageUpdate,
  })
  if err != nil {
    log.Printf("Error when responding the select menu update [%v]", err)
  }
}

// on start fail
func Start_Failed(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {

}

// on start success
// start Game
func Start_Success(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  embed := &dgo.MessageEmbed{
    Title:       "게임 시작!",
    Description: "윙크를 받으셨으면 V 버튼을 클릭 해 주세요!"+
    "\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!"+
    "\n\n**현재 윙크 받은 사람 수 :** 0 / " + strconv.Itoa(guild.Wink.TotalParticipants),
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
    log.Printf("Failed sending follow-up message [%v]", err)
	}
  guild.Wink.MessageIDMap[i.GuildID] = msg.ID
}

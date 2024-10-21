package wink

import (
  // system packages
	"log"

  // internal packages
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// Start
var start_selectMenu dgo.SelectMenu = dgo.SelectMenu{
  CustomID:    "wink_select_list_update",
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

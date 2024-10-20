package wink

import (
  // system packages
	"log"

  // internal packages
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

func SelectListUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  // update currently the currently selected user list
  guild.Wink.SelectedUsersMap[event.GuildID] = event.MessageComponentData().Values

  err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
    // 상호작용 업데이트가 약간의 지연 이후 진행
    Type: dgo.InteractionResponseDeferredMessageUpdate,
  })

  // check if there is an error with the response
  if err != nil {
    log.Printf("Error when responding the select menu update [%v]", err)
  }
}

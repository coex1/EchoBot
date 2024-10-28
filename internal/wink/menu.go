package wink

import (
	// system packages
	"log"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

// on interaction event 'wink_Start_listUpdate'
func Start_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  // update selected user list
  guild.Wink.SelectedUsersID = event.MessageComponentData().Values

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

// on interaction event 'wink_Game_listUpdate'
func Game_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  guild.Wink.UserSelection[event.User.GlobalName] = event.MessageComponentData().Values[0]
}

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
// update selected user list
func Start_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  guild.Wink.SelectedUsersID = event.MessageComponentData().Values
}

// on interaction event 'wink_Start_Button'
func Start_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  players := guild.Wink.SelectedUsersID

  // check if player count is valid TODO: remove comment below
  if len(players) < 2/*MIN_PLAYER_CNT*/ {
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


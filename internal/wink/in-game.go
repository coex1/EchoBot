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

// on interaction event 'wink_norm_list'
func Game_listUpdate(i *dgo.InteractionCreate, g *data.Guild) {
  if g.Wink.State != data.IN_PROGRESS && g.Wink.State != data.LAST_PLAYER {
    return
  }

  g.Wink.UserSelection[i.User.ID] = i.MessageComponentData().Values[0]
}

// when players select their target
func Game_submitButton(s *dgo.Session, i *dgo.InteractionCreate, g *data.Guild) {
  if g.Wink.State != data.IN_PROGRESS && g.Wink.State != data.LAST_PLAYER {
    return
  }

  var target string = g.Wink.UserSelection[i.User.ID]

  g.Wink.UserSelectionFinal[i.User.ID] = target
  log.Printf("User [%s] has selected user [%s] as their target", i.User.GlobalName, g.NameList[target])

	if g.Wink.ConfirmedUsers[i.User.ID] == true {
    return
  }

  // increase total confirmed target count
  g.Wink.ConfirmedUsers[i.User.ID] = true
  g.Wink.ConfirmedCount++
  log.Printf("Confirmed user count [%d/%d]", g.Wink.ConfirmedCount, g.Wink.TotalParticipants)

  // if end conditions are not met, broadcast game status to players
  go func() {
    if end_checkEndCondition(s, g) == false {
      game_broadcastGameStatus(s, g, i.User.GlobalName)
    }
  }()
}

// when the king presses their fake button
func Game_submitKingButton(s *dgo.Session, i *dgo.InteractionCreate, g *data.Guild) {
  if g.Wink.State != data.IN_PROGRESS && g.Wink.State != data.LAST_PLAYER {
    return
  }

  log.Printf("King [%s] has pressed the fake button", i.User.GlobalName)

	if g.Wink.ConfirmedUsers[i.User.ID] == true {
    return
  }

  // increase total confirmed target count
  g.Wink.ConfirmedUsers[i.User.ID] = true
  g.Wink.ConfirmedCount++
  log.Printf("Confirmed user count [%d/%d]", g.Wink.ConfirmedCount, g.Wink.TotalParticipants)

  // if end conditions are not met, broadcast game status to players
  go func() {
    if end_checkEndCondition(s, g) == false {
      game_broadcastGameStatus(s, g, i.User.GlobalName)
    }
  }()
}

// send everyone a message about current game status
// (i.e. who is left to receive a wink)
func game_broadcastGameStatus(s *dgo.Session, guild *data.Guild, u string) {
  var players []string = guild.Wink.SelectedUsersID
  var text_voted string = ""
  var text_not string = ""
  var text string = "현시점 투표 상태!\n\n"
  
  for u, v := range guild.Wink.ConfirmedUsers {
    if v == true {
      text_voted += " -> " + guild.Wink.NameList[u] + "\n"
    } else {
      text_not += " -> " + guild.Wink.NameList[u] + "\n"
    }
  }

  text += "투포 한 사람:\n" + text_voted
  text += "\n투포 안한 사람:\n" + text_not

  message := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "[ " + u + " ]님이 투표하셨습니다!",
        Description:  text,
        Color:        0x26D16D,
      },
    },
  }

  for _, i := range players {
    general.SendComplexDM(s, i, &message)
  }
}

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
  g.Wink.UserSelection[i.User.ID] = i.MessageComponentData().Values[0]
}

// when players select their target
func Game_submitButton(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  var target string = guild.Wink.UserSelection[i.User.ID]

  guild.Wink.UserSelectionFinal[i.User.ID] = target
  log.Printf("User [%s] has selected user [%s] as their target", i.User.GlobalName, guild.NameList[target])

	if guild.Wink.ConfirmedUsers[i.User.ID] == true {
    return
  }

  // increase total confirmed target count
  guild.Wink.ConfirmedUsers[i.User.ID] = true
  guild.Wink.ConfirmedCount++
  log.Printf("Confirmed user count [%d/%d]", guild.Wink.ConfirmedCount, guild.Wink.TotalParticipants)

  // if end conditions are not met, broadcast game status to players
  go func() {
    if checkEndCondition(s, guild) == false {
      broadcastGameStatus(s, guild, i.User.GlobalName)
    }
  }()
}

// when the king presses their fake button
func Game_submitKingButton(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  log.Printf("King [%s] has pressed the fake button", i.User.GlobalName)

	if guild.Wink.ConfirmedUsers[i.User.ID] == true {
    return
  }

  // increase total confirmed target count
  guild.Wink.ConfirmedUsers[i.User.ID] = true
  guild.Wink.ConfirmedCount++
  log.Printf("Confirmed user count [%d/%d]", guild.Wink.ConfirmedCount, guild.Wink.TotalParticipants)

  // if end conditions are not met, broadcast game status to players
  go func() {
    if checkEndCondition(s, guild) == false {
      broadcastGameStatus(s, guild, i.User.GlobalName)
    }
  }()
}

// send everyone a message about current game status
// (i.e. who is left to receive a wink)
func broadcastGameStatus(s *dgo.Session, guild *data.Guild, u string) {
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

  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "[ " + u + " ]님이 투표하셨습니다!",
        Description:  text,
        Color:        0x26D16D,
      },
    },
  }

  for _, i := range players {
    general.SendComplexDM(s, i, &data)
  }
}

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
func Game_listUpdate(e *dgo.InteractionCreate, g *data.Guild) {
  g.Wink.UserSelection[e.User.GlobalName] = e.MessageComponentData().Values[0]
}

func Game_submitButton(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  target := guild.Wink.UserSelection[i.User.GlobalName]
  guild.Wink.UserSelectionFinal[i.User.GlobalName] = target

  log.Printf("User [%s] has selected user [%s] as their target", i.User.GlobalName, target)

	if guild.Wink.ConfirmedUsers[i.User.GlobalName] == true {
    return
  }

  guild.Wink.ConfirmedUsers[i.User.GlobalName] = true
  guild.Wink.ConfirmedCount++
  log.Printf("Confirmed user count [%d]", guild.Wink.ConfirmedCount)

  go func() {
    if checkEndCondition(s, guild) == false {
      broadcastGameStatus(s, guild, i.User.GlobalName)
    }
  }()
}

func Game_submitFakeButton(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  log.Printf("King [%s] has pressed the fake button", i.User.GlobalName)

	if guild.Wink.ConfirmedUsers[i.User.GlobalName] == true {
    return
  }

  guild.Wink.ConfirmedUsers[i.User.GlobalName] = true
  guild.Wink.ConfirmedCount++
  log.Printf("Confirmed user count [%d]", guild.Wink.ConfirmedCount)

  go func() {
    if checkEndCondition(s, guild) == false {
      broadcastGameStatus(s, guild, i.User.GlobalName)
    }
  }()
}

// send everyone a message about current game status
// i.e who is left to receive a wink
func broadcastGameStatus(s *dgo.Session, guild *data.Guild, u string) {
  players := guild.Wink.SelectedUsersID
  text := "**현시점 투표 상태!**\n\n"
  voted := make([]string, guild.Wink.TotalParticipants)
  not := make([]string, guild.Wink.TotalParticipants)
  
  for u, v := range guild.Wink.ConfirmedUsers {
    if v == true {
      voted = append(voted, u)
    } else {
      not = append(not, u)
    }
  }

  text += "**투포 한 사람:**\n"
  for _, v := range voted {
    if len(v) != 0 {
      text += " -> " + v + "\n"
    }
  }

  text += "\n**투포 안한 사람:**\n"
  for _, v := range not {
    if len(v) != 0 {
      text += " -> " + v + "\n"
    }
  }

  embed := dgo.MessageEmbed{
    Title:        "["+u+"]님이 투표하셨습니다!",
    Description:  text,
    Color:        0x26D16D,
  }

  // data for villagers
  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      &embed,
    },
  }

  for _, i := range players {
    general.SendComplexDM(s, i, &data)
  }
}

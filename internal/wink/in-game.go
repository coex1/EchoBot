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

// on interaction event 'wink_Game_listUpdate'
func Game_listUpdate(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
  guild.Wink.UserSelection[event.User.GlobalName] = event.MessageComponentData().Values[0]
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

  // input confirmation DM
  general.SendDM(s, i.User.ID, "지목하신 상대는 [" + target + "] 입니다!\n(원하시면 언제든지 수정하실 수 있으십니다!)")

  go func() {
    if checkEndCondition(s, guild) == false {
      broadcastGameStatus()
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

  // input confirmation DM
  general.SendDM(s, i.User.ID, "윙크 받으셨다고 처리되었습니다!")

  go func() {
    if checkEndCondition(s, guild) == false {
      broadcastGameStatus()
    }
  }()
}

// send everyone a message about current game status
// i.e who is left to receive a wink
func broadcastGameStatus() {

}

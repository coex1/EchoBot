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

func Game_submitButton(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
  target := guild.Wink.UserSelection[i.User.GlobalName]
	guild.Wink.ConfirmedUsers[i.User.GlobalName] = true

  log.Printf("[" + i.User.GlobalName +"] selected user [" + target + "]")

  // ignore index
  general.SendDM(s, i.User.ID, "지목하신 상대는 [" + target + "] 입니다!\n(원하시면 언제든지 수정하실 수 있으십니다!)")

  checkEndCondition(s, guild)
}

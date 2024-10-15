package discord

// internal package
import (
  "github.com/coex1/EchoBot/internal/wink"
	"github.com/coex1/EchoBot/internal/mafia"
)

// external packages
import (
	dgo "github.com/bwmarrin/discordgo"
)

func handleApplicationCommand(s *dgo.Session, event *dgo.InteractionCreate) {
  switch event.ApplicationCommandData().Name {
  case "wink":
    wink.CommandHandle(s, event)
  case "mafia":
    mafia.CommandHandle(s, event)
  }
}



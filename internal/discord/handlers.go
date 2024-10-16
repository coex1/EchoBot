package discord

// system packages
import (
	"log"
)

// internal package
import (
  "github.com/coex1/EchoBot/internal/wink"
	"github.com/coex1/EchoBot/internal/mafia"
	"github.com/coex1/EchoBot/internal/data"
)

// external package
import (
	dgo "github.com/bwmarrin/discordgo"
)

// InteractionCreate event handler
// handler for all user interactions (even Commands!)

// register handlers to 's' session variable
func RegisterHandlers(s *dgo.Session, guild *data.Guild) {
	log.Println("Registering event handlers...")

  s.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
    log.Printf("Login successful [%v#%v]", s.State.User.Username, s.State.User.Discriminator)
  })

  s.AddHandler(func(s *dgo.Session, event *dgo.InteractionCreate) {
    switch event.Type {
    case dgo.InteractionApplicationCommand:
      switch event.ApplicationCommandData().Name {
      case "wink":
        wink.CommandHandle(s, event, guild)
      case "mafia":
        mafia.CommandHandle(s, event)
      }
    case dgo.InteractionMessageComponent:
      switch event.MessageComponentData().CustomID {
      case "wink_select_list_update":
        wink.SelectListUpdate(s, event, guild)
      case "wink_start_button":
        wink.StartButton(s, event, guild)
      case "wink_check", "wink_cancel":
        wink.FollowUpHandler(s, event, guild)
      case "mafia_select_menu":
        mafia.SelectMenu(s, event)
      case "mafia_start_button":
        mafia.StartButton(s, event)
      }
    }
  })
}

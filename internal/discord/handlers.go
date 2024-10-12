package discord

// system packages
import (
	"log"
)

// internal imports
import (
	"github.com/coex1/EchoBot/internal/wink"
)

// external packages
import (
	dgo "github.com/bwmarrin/discordgo"
)

// Ready event handler
// handler for when logged into Discord Server via the Bot Token
var readyEvent = func(s *dgo.Session, r *dgo.Ready) {
	log.Printf("Login successful [%v#%v]", s.State.User.Username, s.State.User.Discriminator)
}

// InteractionCreate event handler
// handler for all user interactions (even Commands!)
var interactionCreateEvent = func(s *dgo.Session, event *dgo.InteractionCreate) {
  switch event.Type {
  case dgo.InteractionApplicationCommand:
    switch event.ApplicationCommandData().Name {
    case "wink":
      wink.SelectUserHandler(s, event)
    case "mafia":
      mafia.SelectUserHandler(s, event)
    }
  case dgo.InteractionMessageComponent:
    switch event.ApplicationCommandData().Name {
    case "wink":
      switch event.MessageComponentData().CustomID {
      case "user_select_menu":
        wink.HandleSelectMenu(s, event) // handleSelectMenu
      case "start_button":
        wink.HandleStartButton(s, event) // winkStartButton
      case "check", "cancel":
        wink.FollowUpHandler(s, event) // winkFollowUpHandler
      }
    case "mafia":
      switch event.MessageComponentData().CustomID {
      case "user_select_menu":
        mafia.handleSelectMenu(s, event)
      case "start_button":
        mafia.mafiaStartButton(s, event) // mafiaStartButton
      }
    }
  }
}

// register handlers to 's' session variable
func RegisterHandlers(s *dgo.Session) {
	log.Println("Registering event handlers...")
	s.AddHandler(readyEvent)
	s.AddHandler(interactionCreateEvent)
}

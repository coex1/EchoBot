package discord

// system packages
import (
	"log"
)

// external package
import (
	dgo "github.com/bwmarrin/discordgo"
)

var (
	winkSelectedUsersMap = make(map[string][]string)
	mafiaSelectedUsersMap = make(map[string][]string)
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
    handleApplicationCommand(s, event)
  case dgo.InteractionMessageComponent:
    switch event.MessageComponentData().CustomID {
    case "wink_user_select_menu":
      Wink_HandleSelectMenu(s, event)
    case "wink_start_button":
      Wink_HandleStartButton(s, event)
    case "wink_check", "wink_cancel":
      Wink_FollowUpHandler(s, event)
    case "mafia_user_select_menu":
      Mafia_HandleSelectMenu(s, event)
    case "mafia_start_button":
      Mafia_HandleStartButton(s, event)
    }
  }
}

// register handlers to 's' session variable
func RegisterHandlers(s *dgo.Session) {
	log.Println("Registering event handlers...")
	s.AddHandler(readyEvent)
	s.AddHandler(interactionCreateEvent)
}

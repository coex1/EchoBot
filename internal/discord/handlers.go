package discord

// system packages
import (
	"log"
)

// external packages
import (
	dgo "github.com/bwmarrin/discordgo"
)

var (
	winkSelectedUsersMap = make(map[string][]string)
	MafiaSelectedUsersMap = make(map[string][]string)
	MinValues        int
	MaxValues        int
)

func handleSelectMenu(s *dgo.Session, event *dgo.InteractionCreate) {
	// Map 변수
  // get currently selected users, and put values to selectedUsersMap
	winkSelectedUsersMap[event.GuildID] = event.MessageComponentData().Values

	err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
		// 상호작용 지연
		Type: dgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to select menu interaction:", err)
	}
}

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
    switch event.ApplicationCommandData().Name {
    case "wink":
      switch event.MessageComponentData().CustomID {
      case "user_select_menu":
        handleSelectMenu(s, event)
      case "start_button":
        wink.HandleStartButton(s, event) // winkStartButton
      case "check", "cancel":
        wink.FollowUpHandler(s, event) // winkFollowUpHandler
      }
    case "mafia":
      switch event.MessageComponentData().CustomID {
      case "user_select_menu":
        handleSelectMenu(s, event)
      case "start_button":
        mafia.StartButton(s, event)
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

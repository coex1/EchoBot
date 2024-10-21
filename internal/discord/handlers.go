package discord

import (
  // system packages
	"log"

  // internal package
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/mafia"
	"github.com/coex1/EchoBot/internal/wink"

  // external package
	dgo "github.com/bwmarrin/discordgo"
)

// register handlers to 's' session variable
func RegisterHandlers(s *dgo.Session, guild *data.Guild) {
	log.Println("Registering event handlers...")

  // Ready(login) event handler
	s.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		log.Printf("Login successful [%v#%v]", s.State.User.Username, s.State.User.Discriminator)
	})

  // InteractionCreate event handler
  // handler for all user interactions (even Commands!)
	s.AddHandler(func(s *dgo.Session, event *dgo.InteractionCreate) {
		switch event.Type {
		case dgo.InteractionApplicationCommand:
			switch event.ApplicationCommandData().Name {
			case "wink":
				wink.CommandHandle(s, event, guild)
			case "mafia":
				mafia.CommandHandle(s, event, guild)
			}
		case dgo.InteractionMessageComponent:
			switch event.MessageComponentData().CustomID {
			case "wink_Start_listUpdate":
				wink.Start_listUpdate(s, event, guild)
			case "wink_Start_Button":
				wink.Start_Button(s, event, guild)

			case "wink_check", "wink_cancel":
				wink.FollowUpHandler(s, event, guild)
			case "mafia_select_menu":
				mafia.SelectMenu(s, event, guild)
			case "mafia_start_button":
				mafia.StartButton(s, event, guild)
			}
		}
	})
}

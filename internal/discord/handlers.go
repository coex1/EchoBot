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
        log.Printf("Starting 'wink' command handle")
				wink.CommandHandle(s, event, guild)
        log.Printf("Finished 'wink' command handle")
			case "mafia":
        log.Printf("Starting 'mafia' command handle")
				mafia.CommandHandle(s, event, guild)
        log.Printf("Finished 'mafia' command handle")
			}
		case dgo.InteractionMessageComponent:
			switch event.MessageComponentData().CustomID {
			case "wink_Start_listUpdate":
        log.Printf("Starting 'wink_Start_listUpdate' handle")
				wink.Start_listUpdate(s, event, guild)
        log.Printf("Finished 'wink_Start_listUpdate' handle")
			case "wink_Start_Button":
        log.Printf("Starting 'wink_Start_Button' handle")
				wink.Start_Button(s, event, guild)
        log.Printf("Finished 'wink_Start_Button' handle")

			case "wink_check":
        log.Printf("Starting 'wink_check' handle")
				wink.FollowUpHandler(s, event, guild)
        log.Printf("Finished 'wink_check' handle")
			case "wink_cancel":
        log.Printf("Starting 'wink_cancel' handle")
				wink.FollowUpHandler(s, event, guild)
        log.Printf("Finished 'wink_cancel' handle")

			case "mafia_select_menu":
				mafia.SelectMenu(s, event, guild)
			case "mafia_start_button":
				mafia.StartButton(s, event, guild)
			}
		}
	})
}

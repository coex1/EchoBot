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
      commandName := event.ApplicationCommandData().Name

      log.Printf("Starting '%s' command handle", commandName)

			switch commandName {
			case "wink":
				wink.Init_CommandHandle(s, event, guild)
			case "mafia":
				mafia.CommandHandle(s, event, guild)
			}

      log.Printf("Finished '%s' command handle", commandName)
		case dgo.InteractionMessageComponent:
      customID := event.MessageComponentData().CustomID

      log.Printf("Starting '%s' handle", customID)

			switch customID {
			case "wink_init_list":
        // send response that event has been received and was acknowledged
        err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
          // Acknowledge that the event has been received,
          // and will be updating the previous message later
          Type: dgo.InteractionResponseDeferredMessageUpdate,
        })
        if err != nil {
          log.Printf("Response to interaction failed [%v]", err)
          return
        }

        wink.Init_listUpdate(event, guild)

			case "wink_start_button":
        err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
          Type: dgo.InteractionResponseDeferredMessageUpdate,
        })
        if err != nil {
          log.Printf("Response to interaction failed [%v]", err)
          return
        }

				wink.Start_buttonPressed(s, event, guild)

			case "wink_norm_list":
        err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
          Type: dgo.InteractionResponseDeferredMessageUpdate,
        })
        if err != nil {
          log.Printf("Response to interaction failed [%v]", err)
          return
        }

				wink.Game_listUpdate(event, guild)

			case "wink_submit_norm_button":
        err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
          Type: dgo.InteractionResponseDeferredMessageUpdate,
        })

        if err != nil {
          log.Printf("Response to interaction failed [%v]", err)
          return
        }

				wink.Game_submitButton(s, event, guild)

			case "wink_submit_king_button":
        err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
          Type: dgo.InteractionResponseDeferredMessageUpdate,
        })

        if err != nil {
          log.Printf("Response to interaction failed [%v]", err)
          return
        }

				wink.Game_submitKingButton(s, event, guild)


			case "wink_end":
        err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
          Type: dgo.InteractionResponseDeferredMessageUpdate,
        })
        if err != nil {
          log.Printf("Response to interaction failed [%v]", err)
          return
        }

        wink.End_Game(guild)

			case "wink_restart":
        err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
          Type: dgo.InteractionResponseDeferredMessageUpdate,
        })
        if err != nil {
          log.Printf("Response to interaction failed [%v]", err)
          return
        }

				wink.Start_Game(s, event, guild)

			case "mafia_select_menu":
				mafia.SelectMenu(s, event, guild)
			case "mafia_start_button":
				mafia.StartButton(s, event, guild)
			}

      log.Printf("Finished '%s' handle", customID)
		}
	})
}

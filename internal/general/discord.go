package general

import (
  // system packages
	"log"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// send a DM(Direct Message) to the user
func SendDM(s *dgo.Session, userID string, msg string) {
	channel, e := s.UserChannelCreate(userID)

  // check if creating a channel with the user failed
	if e != nil {
		log.Printf("Failed creating a direct user channel [%v]", e)
		return
	}

  _, e = s.ChannelMessageSend(channel.ID, msg)

  // check if sending a DM failed
  if e != nil {
    log.Printf("Failed sending DM [%v]", e)
    return
  }

  log.Printf("Sent message to user [%s]", userID)
}

// send a complex DM(Direct Message) to the user
func SendComplexDM(s *dgo.Session, userID string, data *dgo.MessageSend) {
	channel, e := s.UserChannelCreate(userID)

  // check if creating a channel with the user failed
	if e != nil {
		log.Printf("Failed creating a direct user channel [%v]", e)
		return
	}

  _, e = s.ChannelMessageSendComplex(channel.ID, data)

  // check if sending a DM failed
  if e != nil {
    log.Printf("Failed sending DM [%v]", e)
    return
  }

  log.Printf("Sent message to user [%s]", userID)
}

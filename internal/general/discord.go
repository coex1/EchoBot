package general

import (
	// system packages
	"log"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

// send a DM(Direct Message) to the user
func SendDM(s *dgo.Session, userID string, msg string) *dgo.Message {
	channel, e := s.UserChannelCreate(userID)

	// check if creating a channel with the user failed
	if e != nil {
		log.Printf("Failed creating a direct user channel [%v]", e)
		return nil
	}

	m, e := s.ChannelMessageSend(channel.ID, msg)

	// check if sending a DM failed
	if e != nil {
		log.Printf("Failed sending DM [%v]", e)
		return nil
	}

	log.Printf("Sent message to user [%s]", userID)

	return m
}

// send a complex DM(Direct Message) to the user
func SendComplexDM(s *dgo.Session, userID string, data *dgo.MessageSend) *dgo.Message {
	channel, e := s.UserChannelCreate(userID)

	// check if creating a channel with the user failed
	if e != nil {
		log.Printf("Failed creating a direct user channel [%v]", e)
		return nil
	}

	m, e := s.ChannelMessageSendComplex(channel.ID, data)

	// check if sending a DM failed
	if e != nil {
		log.Printf("Failed sending DM [%v]", e)
		return nil
	}

	log.Printf("Sent message to user [%s]", userID)

	return m
}

// send a complex DM(Direct Message) to the user and return both the message and channel ID
func Mafia_SendComplexDM(s *dgo.Session, userID string, data *dgo.MessageSend) (string, *dgo.Message) {
	channel, e := s.UserChannelCreate(userID)

	// check if creating a channel with the user failed
	if e != nil {
		log.Printf("Failed creating a direct user channel [%v]", e)
		return "", nil
	}

	m, e := s.ChannelMessageSendComplex(channel.ID, data)

	// check if sending a DM failed
	if e != nil {
		log.Printf("Failed sending DM [%v]", e)
		return channel.ID, nil
	}

	log.Printf("Sent message to user [%s] in channel [%s]", userID, channel.ID)

	return channel.ID, m
}

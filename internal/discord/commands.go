package discord

import (
	// system packages
	"log"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

// list of commands this bot will register
var cmdList = []*dgo.ApplicationCommand{
	{
		Name:        "wink",
		Description: "윙크게임 시작하기!",
	},
	{
		Name:        "mafia",
		Description: "마피아 게임 시작하기!",
		Options: []*dgo.ApplicationCommandOption{
			{
				Name:        "마피아",
				Description: "Number of Mafia players",
				Type:        dgo.ApplicationCommandOptionInteger,
				Required:    true,
			},
			{
				Name:        "경찰",
				Description: "Number of Police players",
				Type:        dgo.ApplicationCommandOptionInteger,
				MaxValue:    1,
				Required:    true,
			},
			{
				Name:        "의사",
				Description: "Number of Doctor players",
				Type:        dgo.ApplicationCommandOptionInteger,
				MaxValue:    1,
				Required:    true,
			},
		},
	},
	{
		Name:        "echo",
		Description: "test command",
	},
	{
		Name:        "box",
		Description: "test command",
	},
}

// list of commands that successfully registered
var regCmdList []*dgo.ApplicationCommand

// TODO: work without guildID info
func RegisterCommands(s *dgo.Session, guildID string) {
	log.Println("Registering commands...")

	regCmdList = make([]*dgo.ApplicationCommand, len(cmdList))

	for i, v := range cmdList {
		c, e := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)

		// if failed to create command, panic
		if e != nil {
			log.Panicf("Fail to create command '%v' [%v]", v.Name, e)
		}

		regCmdList[i] = c
	}

	log.Println("Successfully registered commands!")
}

// TODO: work without guildID info
func RemoveCommands(s *dgo.Session, guildID string) {
	log.Println("Removing commands...")

	for _, v := range regCmdList {
		e := s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)

		// if failed to create command, panic
		if e != nil {
			log.Panicf("Failed to remove command '%v' [%v]", v.Name, e)
		}
	}

	log.Println("Successfully removed commands!")
}

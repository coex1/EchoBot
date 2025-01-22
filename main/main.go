package main

import (
	// system packages
	"log"

	// internal imports
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/discord"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

// bot configuration values read from file
var config Config

// global running-data holding struct variable
var guild *data.Guild

// all information regarding the connection to Discord Server
var discordSession *dgo.Session

// initialize a session Discord servers
func init() {
	var e error

	// get configuration value
	GetBotConfiguration(&config)

	// create running data
	guild = new(data.Guild) // TODO: need to be data array, to hold multiple guilds
	data.Initialize(guild)

	// create a new session, and initialize it with the BotTokenKey
	discordSession, e = dgo.New("Bot " + config.BotTokenKey)

	// if error is found
	if e != nil {
		log.Fatalf("invalid bot parameters [%v]", e)
	}
}

func main() {
	discord.RegisterHandlers(discordSession, guild)

	discord.Start(discordSession)
	defer discord.Stop(discordSession)

	discord.RegisterCommands(discordSession, config.GuildID)

	WaitOnInterrupt()

	discord.RemoveCommands(discordSession, config.GuildID)
}

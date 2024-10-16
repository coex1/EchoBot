package main

// system packages
import (
	"flag"
	"log"
)

// internal imports
import (
	"github.com/coex1/EchoBot/internal/discord"
	"github.com/coex1/EchoBot/internal/data"
)

// external packages
import (
	dgo "github.com/bwmarrin/discordgo"
)

// TODO: remove
var (
	BotToken = flag.String("token", "", "")
	GuildID  = flag.String("guild", "948807733199642645", "")
	//RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
  guild = new(data.Guild) // TODO: need to be data array
)

// TODO: remove
func init() {
	flag.Parse()
}

// all information regarding the connection to Discord Server
var discordSession *dgo.Session

// initialize a session Discord servers
func init() {
	var e error

	// TODO: change BotToken type to just string, not string pointer
	discordSession, e = dgo.New("Bot " + *BotToken)

	// if error is found
	if e != nil {
		log.Fatalf("invalid bot parameters [%v]", e)
	}
}

func main() {
	discord.RegisterHandlers(discordSession, guild)

	discord.Start(discordSession)
	defer discord.Stop(discordSession)

	discord.RegisterCommands(discordSession, *GuildID)

	WaitOnInterrupt()

	discord.RemoveCommands(discordSession, *GuildID)
}

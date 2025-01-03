package wink

import (
  // system packages
  "log"

  // internal packages
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// reset all global variables
func reset_fully(g *data.Guild) {
  log.Printf("Resetting game data (fully)\n")

  // game state
  g.Wink.State = 0

  // a list of the channel's members (before selection)
	g.Wink.AllUserInfo = make([]dgo.SelectMenuOption, 0)

	g.Wink.SelectedUsersID = make([]string, 0)
	g.Wink.SelectedUsersInfo = make([]dgo.SelectMenuOption, 0)
	g.Wink.TotalParticipants = 0

	g.Wink.ConfirmedUsers = make(map[string]bool)

  // ?
	g.Wink.CheckedUsers = make(map[string]bool)
	g.Wink.MessageIDMap = make(map[string]string)
  g.Wink.UserSelection = make(map[string]string)
  g.Wink.UserSelectionFinal = make(map[string]string)
}

// reset all session variables
func reset_part(g *data.Guild) {
  log.Printf("Resetting game data (partially)\n")

	g.Wink.ConfirmedUsers = make(map[string]bool)

  // ?
	g.Wink.CheckedUsers = make(map[string]bool)
	g.Wink.MessageIDMap = make(map[string]string)
  g.Wink.UserSelection = make(map[string]string)
}

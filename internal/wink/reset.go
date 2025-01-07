package wink

import (
  // system packages
  "log"

  // internal packages
  "github.com/coex1/EchoBot/internal/data"
)

// reset all global variables
func resetGame(g *data.Guild) {
  log.Printf("Resetting game data (fully)\n")

  // game state
  g.Wink.State = 0

  // a list of the channel's members (before selection)

	g.Wink.SelectedUsersID = make([]string, 0)
	g.Wink.TotalParticipants = 0

  g.Wink.MaxPossiblePlayers = 0

	g.Wink.ConfirmedUsers = make(map[string]bool)
  g.Wink.ConfirmedCount = 0
  g.Wink.KingID = "" 
  g.Wink.KingName = "" 

	g.Wink.NameList = make(map[string]string)
	g.Wink.IDList = make(map[string]string)

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
  g.Wink.ConfirmedCount = 0

  // ?
	g.Wink.CheckedUsers = make(map[string]bool)
	g.Wink.MessageIDMap = make(map[string]string)
  g.Wink.UserSelection = make(map[string]string)
}

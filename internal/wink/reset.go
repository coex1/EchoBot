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

  g.Wink.State = 0

	g.Wink.NameList = make(map[string]string)
  g.Wink.MaxPossiblePlayers = 0
  g.Wink.SelectedUsersID = make([]string, 0)
  g.Wink.TotalParticipants = 0

  resetPart(g)
}

// reset all session variables
func resetPart(g *data.Guild) {
  log.Printf("Resetting session data\n")

	g.Wink.ConfirmedUsers = make(map[string]bool)
  g.Wink.ConfirmedCount = 0

	g.Wink.ConfirmedUsers = make(map[string]bool)
  g.Wink.ConfirmedCount = 0

  g.Wink.KingID = "" 
  g.Wink.FinalPlayerID = ""

  g.Wink.UserSelection = make(map[string]string)
  g.Wink.UserSelectionFinal = make(map[string]string)
}

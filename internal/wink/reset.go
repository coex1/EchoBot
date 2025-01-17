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
func resetGame(g *data.Guild) {
  log.Printf("Resetting game data (fully)\n")

  g.Wink.State = 0

	g.Wink.NameList = make(map[string]string)
  g.Wink.MaxPossiblePlayers = 0
  g.Wink.SelectedUsersID = make([]string, 0)
  g.Wink.TotalParticipants = 0

  g.Wink.TargetList = make([]dgo.SelectMenuOption, 0)  // a menu list of players to select

  resetPart(g)
}

// reset all session variables
func resetPart(g *data.Guild) {
  log.Printf("Resetting session data\n")

	g.Wink.ConfirmedUsers = make(map[string]bool)
  g.Wink.ConfirmedCount = 0

  g.Wink.KingID = "" 
  g.Wink.FinalPlayerID = ""

  g.Wink.UserSelection = make(map[string]string)
  g.Wink.UserSelectionFinal = make(map[string]string)
}

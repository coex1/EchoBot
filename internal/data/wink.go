package data

import (
  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// game states
const (
  NONE = iota + 0
  INITIATED
  IN_PROGRESS
  LAST_PLAYER
  ENDED
)

type Wink struct {
  State               int                     // game state

  NameList            map[string]string       // list of every possible player
  MaxPossiblePlayers  int
	SelectedUsersID     []string                // users selected to play the game
	TotalParticipants   int

  TargetList          []dgo.SelectMenuOption  // a menu list of players to select

	ConfirmedUsers      map[string]bool         // users that have confirmed their target (ID -> BOOL)
  ConfirmedCount      int

  KingID              string                  // king's ID
  FinalPlayerID       string                  // final player

	UserSelection       map[string]string       // player's selection (ID -> ID)
	UserSelectionFinal  map[string]string       // player's final selection (ID -> ID)
}

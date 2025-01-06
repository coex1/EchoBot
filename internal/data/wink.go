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
  // game state
  State             int

  // all user detail
	AllUserInfo       []dgo.SelectMenuOption

  // all selected detail
	SelectedUsersInfo []dgo.SelectMenuOption

  // users selected to play the game
	SelectedUsersID   []string

  // users that have confirmed their target
	ConfirmedUsers    map[string]bool
  ConfirmedCount    int

  // kingID
  KingID string
  KingName string

  MasterList map[string]string

	CheckedUsers      map[string]bool
	TotalParticipants int
	MessageIDMap      map[string]string

	UserSelection     map[string]string
	UserSelectionFinal     map[string]string
}

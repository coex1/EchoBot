package data

import (
  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

type Wink struct {
  // all user detail
	AllUserInfo       []dgo.SelectMenuOption

  // all selected detail
	SelectedUsersInfo []dgo.SelectMenuOption

  // users selected to play the game
	SelectedUsersID   []string

  // users that have confirmed their target
	ConfirmedUsers    map[string]bool

	CheckedUsers      map[string]bool
	TotalParticipants int
	MessageIDMap      map[string]string

	UserSelection     map[string]string
}

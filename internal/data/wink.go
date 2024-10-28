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

	CheckedUsers      map[string]bool
	TotalParticipants int
	MessageIDMap      map[string]string

	UserSelection     map[string]string
}

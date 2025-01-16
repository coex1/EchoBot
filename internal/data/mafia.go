package data

import dgo "github.com/bwmarrin/discordgo"

type Mafia struct {
	MasterList map[string]string

	// all user detail
	AllUserInfo []dgo.SelectMenuOption

	// users selected to play the game
	SelectedUsersID []string

	MafiaList   []string
	PoliceList  []string
	DoctorList  []string
	CitizenList []string
}

package data

import dgo "github.com/bwmarrin/discordgo"

type Mafia struct {
	NumMafia  int
	NumPolice int
	NumDoctor int

	// day : 0 | night : 1
	State int

	// all user detail
	AllUserInfo []dgo.SelectMenuOption

	// users selected to play the game
	SelectedUsersID []string

	MafiaList   []string
	PoliceList  []string
	DoctorList  []string
	CitizenList []string
}

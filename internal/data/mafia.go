package data

import (
	dgo "github.com/bwmarrin/discordgo"
)

type Mafia struct {
	// both
	SelectedUsersID []string // ID

	// start
	NumMafia  int
	NumPolice int
	NumDoctor int

	AllUserInfo []dgo.SelectMenuOption // 이름 : ID

	UserDMChannels map[string]string

	// in-game
	Timer int

	State bool // day : 1 | night : 0

	VoteMap   map[string]string // ID : Vote_ID
	VoteCount map[string]int    // ID : Count

	AliveUsersID []string

	MafiaList   []string
	PoliceList  []string
	DoctorList  []string
	CitizenList []string
}

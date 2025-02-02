package data

import (
	dgo "github.com/bwmarrin/discordgo"
)

type Mafia struct {
	// both
	SelectedUsersID []string               // ID
	AllUserInfo     []dgo.SelectMenuOption // 이름 : ID

	// start
	NumMafia  int
	NumPolice int
	NumDoctor int

	ReadyMap map[string]bool

	// in-game
	Timer int

	Day int

	State bool // day : 1 | night : 0

	VoteMap   map[string]string // ID : Vote_ID
	VoteCount map[string]int    // ID : Count

	AliveUsersID []dgo.SelectMenuOption

	MafiaList   []string
	PoliceList  []string
	DoctorList  []string
	CitizenList []string
}

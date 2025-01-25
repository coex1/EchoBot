package data

import dgo "github.com/bwmarrin/discordgo"

type Mafia struct {
	// both
	SelectedUsersID []string // ID

	// start
	NumMafia  int
	NumPolice int
	NumDoctor int

	MessageIDMap map[string]string

	AllUserInfo []dgo.SelectMenuOption // 이름 : ID

	// in-game
	Timer int

	State bool // day : 1 | night : 0

	VoteList  map[string]string // ID : Vote_ID
	VoteCount map[string]int    // ID : Vote

	AliveUserInfo []dgo.SelectMenuOption

	MafiaList   []string
	PoliceList  []string
	DoctorList  []string
	CitizenList []string
}

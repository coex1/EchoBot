package data

import (
	dgo "github.com/bwmarrin/discordgo"
)

type Mafia struct {
	// 모든 플레이어 정보 (ID -> 플레이어 정보)
	Players map[string]*MafiaPlayer

	// 선택된 플레이어 목록
	SelectedUsersID []string               // ID
	AllUserInfo     []dgo.SelectMenuOption // 이름 : ID

	// 게임 시작 시 인원 수
	NumMafia  int
	NumPolice int
	NumDoctor int

	// 준비 여부
	ReadyMap map[string]bool

	// 게임 진행 정보
	Timer int

	Day int

	State bool // day : 1 | night : 0

	// 투표 정보
	VoteMap   map[string]string // ID : Vote_ID
	VoteCount map[string]int    // ID : Count

	// 생존 정보
	AliveUsersIDMap []dgo.SelectMenuOption
	AliveUsersID    []string

	// 역할 리스트
	MafiaList   []string
	PoliceList  []string
	DoctorList  []string
	CitizenList []string
}

type MafiaPlayer struct {
	ID          string // Discord User ID
	GlobalName  string // 유저 닉네임
	DMChannelID string // DM 채널 ID
	Role        string // 역할 (Mafia, Police, Doctor, Citizen)
	IsAlive     bool   // 생존 여부
}

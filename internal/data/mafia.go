package data

import (
	"context"

	dgo "github.com/bwmarrin/discordgo"
)

type Mafia struct {
	// 모든 플레이어 정보 (ID -> 플레이어 정보)
	Players   map[string]*MafiaPlayer
	ChannelID string

	// 선택된 플레이어 목록
	SelectedUsersID []string               // ID
	AllUserInfo     []dgo.SelectMenuOption // 이름 : ID

	// 게임 시작 시 인원 수
	NumMafia  int
	NumPolice int
	NumDoctor int

	SleepPhrases []string

	// 게임 진행 정보
	TimerActive bool

	Day int

	// 투표 정보 (Reset)
	TempVoteMap map[string]string
	VoteMap     map[string]string // ID : Vote_ID
	VoteCount   map[string]int    // ID : Count
	VoteSkip    []string          // ID

	MafiaTargetMap map[string]string // Used Night Action
	MafiaTarget    string
	PoliceTarget   string
	DoctorTarget   string
	CitizenPhrases map[string]string
	CitizenReady   map[string]bool

	CancelFunc context.CancelFunc

	// Day and Night
	State bool

	// 스킬 사용 상태
	NightActionDone map[string]bool // role : bool (Mafia, Police, Doctor, Citizen)
}

type MafiaPlayer struct {
	ID          string
	GlobalName  string
	DMChannelID string
	Role        string // (Mafia, Police, Doctor, Citizen)
	IsAlive     bool
}

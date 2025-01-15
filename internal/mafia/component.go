package mafia

import (
	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

// Start
var min_player_cnt = MIN_PLAYER_CNT
var start_selectMenu dgo.SelectMenu = dgo.SelectMenu{
	CustomID:    "mafia_select_menu",
	Placeholder: "사용자를 선택해 주세요!",
	MinValues:   &min_player_cnt,
	MaxValues:   MAX_MEMBER_GET,
	Options:     []dgo.SelectMenuOption{},
}

// start_button
var start_buttonRow dgo.ActionsRow = dgo.ActionsRow{
	Components: []dgo.MessageComponent{
		&dgo.Button{
			Label:    "게임 시작",              // 버튼 텍스트
			Style:    dgo.PrimaryButton,    // 버튼 스타일
			CustomID: "mafia_start_button", // 버튼 클릭 시 처리할 ID
		},
	},
}

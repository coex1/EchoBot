package mafia

import (
	// system packages
	"log"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

func CommandHandle(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	var optionList []dgo.SelectMenuOption
	var minListCnt int = MAFIA_MIN_LIST_CNT
	var maxListCnt int
	var err error
	var members []*dgo.Member

	guild.Mafia.SelectedUsersMap = make(map[string][]string)

	// get guild members
	members, err = s.GuildMembers(event.GuildID, QUERY_STRING, MAX_MEMBER_GET)

	// if getting members failed
	if err != nil {
		log.Fatalf("Failed getting members [%v]", err)
		return
	}

	// create select list from 'members'
	for _, m := range members {
		// check if 'm' is a bot
		if m.User.Bot {
			continue
		}

		optionList = append(optionList, dgo.SelectMenuOption{
			Label: m.User.GlobalName,
			Value: m.User.ID,
		})
	}

	// update max to match the updated 'optionList'
	maxListCnt = len(optionList)

	// SelectMenu와 ActionRow 설정
	selectMenu := dgo.SelectMenu{
		CustomID:    "mafia_user_select_menu",
		Placeholder: "사용자를 선택해 주세요!",
		MinValues:   &minListCnt,
		MaxValues:   maxListCnt,
		Options:     optionList,
	}
	actionRow := dgo.ActionsRow{
		Components: []dgo.MessageComponent{
			selectMenu,
		},
	}

	// start_button
	buttonRow := dgo.ActionsRow{
		Components: []dgo.MessageComponent{
			&dgo.Button{
				Label:    "시작",                 // 버튼 텍스트
				Style:    dgo.PrimaryButton,    // 버튼 스타일
				CustomID: "mafia_start_button", // 버튼 클릭 시 처리할 ID
			},
		},
	}

	// 드롭다운 메뉴와 버튼을 포함한 메시지 전송
	err = s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Components: []dgo.MessageComponent{
				actionRow,
				buttonRow,
			},
		},
	})

	// if response failed
	if err != nil {
		log.Fatalf("Failed to send response [%v]", err)
		return
	}

	// TODO: change from global variable to local
	// reset selected users mapping
	guild.Mafia.SelectedUsersMap[event.GuildID] = make([]string, 0)
}

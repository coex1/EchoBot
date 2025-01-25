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
	var err error
	var members []*dgo.Member
	var minListCnt int = MIN_PLAYER_CNT

	// get guild members
	members, err = s.GuildMembers(event.GuildID, QUERY_STRING, MAX_MEMBER_GET)
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

		guild.Mafia.AllUserInfo = append(guild.Mafia.AllUserInfo, dgo.SelectMenuOption{
			Label: m.User.GlobalName,
			Value: m.User.ID,
		})
	}

	// 드롭다운 메뉴와 버튼을 포함한 메시지 전송
	response := &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Embeds: []*dgo.MessageEmbed{
				{
					Title: "게임 참여자 선택!",
					Description: "게임에 참석할 사용자들을 선택해 주세요.\n\n" +
						"**최소 5명** 이상이 선택 되어야 게임이 가능합니다.\n\n" +
						"선택 하셨으면 **게임시작** 버튼을 클릭 해 주세요.",
					Color: 0x2AFF00,
				},
			},
			Components: []dgo.MessageComponent{
				dgo.ActionsRow{
					Components: []dgo.MessageComponent{
						dgo.SelectMenu{
							MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
							CustomID:    "mafia_Start_listUpdate",
							Placeholder: "사용자 목록",
							MinValues:   &minListCnt,
							MaxValues:   len(guild.Mafia.AllUserInfo),
							Options:     guild.Mafia.AllUserInfo,
						},
					},
				},
				dgo.ActionsRow{
					Components: []dgo.MessageComponent{
						&dgo.Button{
							CustomID: "mafia_Start_Button", // 버튼 클릭 시 처리할 ID
							Label:    "게임시작",               // 버튼 텍스트
							Style:    dgo.PrimaryButton,    // 버튼 스타일
						},
					},
				},
			},
		},
	}
	err = s.InteractionRespond(event.Interaction, response)
	if err != nil {
		log.Fatalf("failed to send response [%v]", err)
		return
	}
}

package mafia

import (
	// system packages
	"log"

	// internal packages
	"github.com/coex1/EchoBot/internal/data"

	// external packages
	dgo "github.com/bwmarrin/discordgo"
)

func CommandHandle(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	var minListCnt int = MIN_PLAYER_CNT
	var optionList []dgo.SelectMenuOption
	var err error
	var members []*dgo.Member

	guild.Mafia.SelectedUsers = make(map[string][]string)

	// get guild members
	members, err = s.GuildMembers(i.GuildID, QUERY_STRING, MAX_MEMBER_GET)
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

	start_selectMenu.MinValues = &minListCnt
	start_selectMenu.MaxValues = len(optionList)
	start_selectMenu.Options = optionList

	// 드롭다운 메뉴와 버튼을 포함한 메시지 전송
	err = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Components: []dgo.MessageComponent{
				dgo.ActionsRow{
					Components: []dgo.MessageComponent{
						start_selectMenu,
					},
				},
				start_buttonRow,
			},
		},
	})
	// if response failed
	if err != nil {
		log.Fatalf("Failed to send response [%v]", err)
		return
	}
}

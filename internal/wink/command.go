package wink

import (
  // system packages
  "log"

  // internal packages
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)


func CommandHandle(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	var minListCnt int = MIN_PLAYER_CNT
	var err error
	var optionList []dgo.SelectMenuOption
	var members []*dgo.Member

	guild.Wink.SelectedUsers[event.GuildID] = make([]string, 0)

	guild.Wink.CheckedUsers = make(map[string]bool)
	guild.Wink.TotalParticipants = 0
	guild.Wink.MessageIDMap = make(map[string]string)

  log.Printf("Starting command handler\n")

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

		optionList = append(optionList, dgo.SelectMenuOption{
			Label: m.User.GlobalName,
			Value: m.User.ID,
		})
	}

  start_selectMenu.MinValues = &minListCnt
  start_selectMenu.MaxValues = len(optionList)
  start_selectMenu.Options = optionList

  // respond to command by sending Start Menu
	err = s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
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
	if err != nil {
		log.Fatalf("Failed to send response [%v]", err)
		return
	}

  log.Printf("Finished command handler\n")
}

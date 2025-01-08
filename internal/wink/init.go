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
	var err error
	var members []*dgo.Member // for getting channel members

  if guild.Wink.State != data.NONE && guild.Wink.State != data.ENDED {
		log.Printf("Invalid game state to run this command")
    init_sendFailedResponse(s, event, "현재 상태에서 이 명령어를 사용하실 수 없습니다!")
    return
  }

  // reset all global variables
  resetGame(guild)

	// get guild members
	members, err = s.GuildMembers(event.GuildID, QUERY_STRING, MAX_MEMBER_GET)
	if err != nil {
		log.Printf("Failed getting members [%v]", err)
    init_sendFailedResponse(s, event, "Failed to get members")
		return
	}

	// create select list from 'members'
	for _, m := range members {
		// check if 'm' is a bot
		if m.User.Bot {
			continue
		}

    guild.Wink.NameList[m.User.ID] = m.User.GlobalName
    guild.Wink.IDList[m.User.GlobalName] = m.User.ID
    guild.Wink.MaxPossiblePlayers++
	}

  init_sendGameInitResponse(s, event, guild)
}

func init_sendFailedResponse(s *dgo.Session, e *dgo.InteractionCreate, f string) {
  // create reponse
  response := &dgo.InteractionResponse{
    Type: dgo.InteractionResponseChannelMessageWithSource,
    Data: &dgo.InteractionResponseData{
      Embeds: []*dgo.MessageEmbed{
        {
          Title:        "시작 실패",
          Description:  "아래 이유로 게임을 시작하기 실패하였습니다.\n"+
                        "["+f+"]",
          Color:        0xC71818,
        },
      },
    },
  }

  // send response
  err := s.InteractionRespond(e.Interaction, response)
  if err != nil {
    log.Printf("Failed to send response [%v] (g2)", err)
    return
  }
}

func init_sendGameInitResponse(s *dgo.Session, e *dgo.InteractionCreate, g *data.Guild) {
	var min int = MIN_PLAYER_CNT
  var max int = len(g.NameList)
	var list []dgo.SelectMenuOption = make([]dgo.SelectMenuOption, 0)
	//var list []dgo.SelectMenuOption = make([]dgo.SelectMenuOption, max)

  log.Printf("DEBUG: 1[%d] 2[%d]\n", g.Wink.MaxPossiblePlayers, max)

  // create menu list
  for id, name := range g.NameList {
    list = append(list, dgo.SelectMenuOption{
      Label: name,
      Value: id,
    })
    log.Printf("name: [%s], id: [%s]\n", name, id);
  }

  // create reponse
  response := &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
      Embeds: []*dgo.MessageEmbed{ 
        {
          Title:        "게임 참여자 선택!",
          Description:  "게임에 참석할 사용자들을 선택해 주세요.\n"+
                        "최소 5명 이상이 선택 되어야 게임이 가능합니다.\n"+
                        "선택 하셨으면 '게임시작' 버튼을 클릭 해 주세요.",
          Color:        0x2AFF00,
        },
      },
			Components: []dgo.MessageComponent{
        dgo.ActionsRow{
          Components: []dgo.MessageComponent{
            dgo.SelectMenu{
              MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
              CustomID:    "wink_init_list",
              Placeholder: "사용자 목록",
              MinValues:   &min,
              MaxValues:   len(list),
              Options:     list,
            },
          },
        },
        dgo.ActionsRow{
          Components: []dgo.MessageComponent{
            &dgo.Button{
              CustomID: "wink_start_button", // 버튼 클릭 시 처리할 ID
              Label:    "게임시작",          // 버튼 텍스트
              Style:    dgo.PrimaryButton,   // 버튼 스타일
            },
          },
        },
			},
		},
	}

  // send response
  err := s.InteractionRespond(e.Interaction, response)
	if err != nil {
		log.Printf("Failed to send response [%v] (g1)", err)
		return
	}
}

// on interaction event 'wink_init_list'
// update selected user list
func Init_listUpdate(e *dgo.InteractionCreate, g *data.Guild) {
  g.Wink.SelectedUsersID = e.MessageComponentData().Values
}

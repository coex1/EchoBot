package mafia

// system packages
import (
	"log"
	"math/rand"
	"time"
)

// internal packages
import (
  "github.com/coex1/EchoBot/internal/general"
)

// external package
import (
	dgo "github.com/bwmarrin/discordgo"
)

const MAFIA_MIN_LIST_CNT = 3
const MAX_MEMBER_GET int = 50
const QUERY_STRING string = ""

var (
  // 윙크 받아서 버튼은 클릭 한 사용자들 
	selectedUsersMap = make(map[string][]string)
)

func SelectMenu(s *dgo.Session, event *dgo.InteractionCreate) {
	// Map 변수
  // get currently selected users, and put values to selectedUsersMap
	selectedUsersMap[event.GuildID] = event.MessageComponentData().Values

	err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
		// 상호작용 지연
		Type: dgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to select menu interaction:", err)
	}
}

func StartButton(s *dgo.Session, event *dgo.InteractionCreate) {
	mafiaCount := event.ApplicationCommandData().Options[0].IntValue()
	policeCount := event.ApplicationCommandData().Options[1].IntValue()
	doctorCount := event.ApplicationCommandData().Options[2].IntValue()

	var mafiaSelected, policeSelected, doctorSelected []string
	totalCount := int(mafiaCount + policeCount + doctorCount)

	tempSelectedMembers := selectedUsersMap[event.GuildID]
	if len(tempSelectedMembers) == 0 {
		log.Println("No members selected.")
		return
	}

	if totalCount > len(tempSelectedMembers) {
		log.Println("Exceeded number")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]string, 0)
	copy(shuffled, tempSelectedMembers)
	r.Shuffle(len(shuffled), func(event, j int) {
		shuffled[event], shuffled[j] = shuffled[j], shuffled[event]
	})
	mafiaSelected = shuffled[:mafiaCount]
	shuffled = shuffled[mafiaCount:]
	policeSelected = shuffled[:policeCount]
	shuffled = shuffled[policeCount:]
	doctorSelected = shuffled[:doctorCount]

	var message string
	for _, id := range tempSelectedMembers {
		if general.Contains(mafiaSelected, id) {
			message = "당신은 마피아 입니다!"
		} else if general.Contains(policeSelected, id) {
			message = "당신은 경찰 입니다!"
		} else if general.Contains(doctorSelected, id) {
			message = "당신은 의사 입니다!"
		} else {
			message = "당신은 시민 입니다!"
		}
		general.SendDM(s, id, message)
	}
}

func CommandHandle(s *dgo.Session, event *dgo.InteractionCreate) {
	var optionList []dgo.SelectMenuOption
  var minListCnt int = MAFIA_MIN_LIST_CNT
  var maxListCnt int
  var err error
  var members []*dgo.Member

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
				Label:    "시작",                    // 버튼 텍스트
				Style:    dgo.PrimaryButton, // 버튼 스타일
				CustomID: "mafia_start_button",          // 버튼 클릭 시 처리할 ID
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
	selectedUsersMap[event.GuildID] = make([]string, 0)
}

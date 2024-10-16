package wink

// system packages
import (
	"fmt"
	"log"
	"time"
	"math/rand"
	"strconv"
	"strings"
)

// internal packages
import (
  "github.com/coex1/EchoBot/internal/general"
)

// external packages
import (
	dgo "github.com/bwmarrin/discordgo"
)

const WINK_MIN_LIST_CNT = 2
const MAX_MEMBER_GET int = 50
const QUERY_STRING string = ""

var (
  // 윙크 받아서 버튼은 클릭 한 사용자들 
	checkedUsers      = make(map[string]bool)
	totalParticipants = 0
	messageIDMap      = make(map[string]string)
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

func FollowUpHandler(s *dgo.Session, event *dgo.InteractionCreate) {
	userID := event.Member.User.ID
	userGlobalName := event.Member.User.GlobalName
	action := ""

	// 버튼 클릭에 따라 상태 업데이트
	switch event.MessageComponentData().CustomID {
	case "check":
		if !checkedUsers[userID] {
			checkedUsers[userID] = true
		}
		action = "확인"
	case "cancel":
		if checkedUsers[userID] {
			checkedUsers[userID] = false
		}
		action = "취소"
	}

	// 현재 체크된 수 계산
	checkedCount := general.CountCheckedUsers(checkedUsers)

	// 기존 메시지 업데이트 (Followup 메시지 수정)
	messageID := messageIDMap[event.GuildID]
	if messageID == "" {
		log.Println("No message found to update")
		return
	}

	// 체크된 유저 및 체크되지 않은 유저 목록 생성
	var checkedUsersList, uncheckedUsersList string
	for _, id := range selectedUsersMap[event.GuildID] {
		// 유저 정보를 가져오기
		member, err := s.GuildMember(event.GuildID, id)
		if err != nil {
			log.Println("Error fetching member:", err)
			continue
		}
		userName := member.User.GlobalName

		if checkedUsers[id] {
			checkedUsersList += fmt.Sprintf("%s\n", userName)
		} else {
			uncheckedUsersList += fmt.Sprintf("%s\n", userName)
		}
	}

	// 남은 사람이 한 명일 경우 처리
	var embed *dgo.MessageEmbed
	if checkedCount == totalParticipants-1 {
		lastUserName := uncheckedUsersList
		embed = &dgo.MessageEmbed{
			Title: "마지막 남은 사람!\n",
			Description: fmt.Sprintf(
				"%s님, 당신이 마지막 사람입니다.\n\n왕일 것 같은 사람을 지목해주세요!", strings.ReplaceAll(lastUserName, "\n", ""),
			),
			Color: 0xff0000, // 다른 색으로 표시
		}
	} else {
		embed = &dgo.MessageEmbed{
			Title: "게임 진행 중...\n",
			Description: fmt.Sprintf(
				"윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** %d / %d\n\n**확인한 유저 :**\n%s\n**확인하지 못한 유저 :**\n%s",
				checkedCount, totalParticipants, checkedUsersList, uncheckedUsersList,
			),
			Color: 0x00ff00,
		}
	}

	// 상호작용 응답 지연 후 아래 메시지 수정 진행
	err := s.InteractionRespond(event.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	// 메시지 수정
	content := fmt.Sprintf("'%s'이(가) %s했습니다.\n", userGlobalName, action)
	_, err = s.ChannelMessageEditComplex(&dgo.MessageEdit{
		Channel:    event.ChannelID,
		ID:         messageID,
		Embeds:     &[]*dgo.MessageEmbed{embed},
		Content:    &content,
		Components: &event.Message.Components, // 기존 버튼 컴포넌트 유지
	})
	if err != nil {
		log.Println("Error updating message:", err)
		return
	}
}

func FollowUpMessage(s *dgo.Session, i *dgo.InteractionCreate) {
	totalParticipants = len(selectedUsersMap[i.GuildID])

	embed := &dgo.MessageEmbed{
		Title:       "게임 시작!",
		Description: "윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** 0 / " + strconv.Itoa(totalParticipants),
		Color:       0x00ff00,
	}

	components := []dgo.MessageComponent{
		dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				&dgo.Button{
					Label:    "V",
					Style:    dgo.SuccessButton,
					CustomID: "wink_check",
				},
				&dgo.Button{
					Label:    "X",
					Style:    dgo.DangerButton,
					CustomID: "wink_cancel",
				},
			},
		},
	}

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &dgo.WebhookParams{
		Embeds:     []*dgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		log.Println("Error sending follow-up message:", err)
	}

	messageIDMap[i.GuildID] = msg.ID
}

// 버튼 포함 임베드 메시지 생성
func CreateFollowUpMessage(s *dgo.Session, i *dgo.InteractionCreate) {
	totalParticipants = len(selectedUsersMap[i.GuildID])

	embed := &dgo.MessageEmbed{
		Title:       "게임 시작!",
		Description: "윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** 0 / " + strconv.Itoa(totalParticipants),
		Color:       0x00ff00,
	}

	components := []dgo.MessageComponent{
		dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				&dgo.Button{
					Label:    "V",
					Style:    dgo.SuccessButton,
					CustomID: "wink_check",
				},
				&dgo.Button{
					Label:    "X",
					Style:    dgo.DangerButton,
					CustomID: "wink_cancel",
				},
			},
		},
	}

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &dgo.WebhookParams{
		Embeds:     []*dgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		log.Println("Error sending follow-up message:", err)
	}

	messageIDMap[i.GuildID] = msg.ID
}

// Select 버튼이 눌렸을 때 선택된 멤버들을 처리하는 핸들러
func StartButton(s *dgo.Session, i *dgo.InteractionCreate) {
	// 선택된 멤버 ID 목록을 가져옴
	tempSelectedMembers := selectedUsersMap[i.GuildID]
	if len(tempSelectedMembers) == 0 {
		log.Println("No members selected.")
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate a random number between 0 and 100
	randomNumber := r.Intn(len(tempSelectedMembers)) // Intn(n) returns a random integer from 0 to n-1, so 101 gives 0 to 100
	fmt.Println(randomNumber)                        // Print the random number

	king := tempSelectedMembers[randomNumber]

	var message string
	for _, id := range tempSelectedMembers {
		if id == king {
			message = "당신은 왕 입니다!"
		} else {
			message = "당신은 왕이 아닙니다!"
		}
	  general.SendDM(s, id, message)
	}

	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	FollowUpMessage(s, i)
}

func CommandHandle(s *dgo.Session, event *dgo.InteractionCreate) {
	var optionList []dgo.SelectMenuOption
  var minListCnt int = WINK_MIN_LIST_CNT
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

  // configure select menu
	selectMenu := dgo.SelectMenu{
		CustomID:    "wink_user_select_menu",
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

	// configure game start button
	buttonRow := dgo.ActionsRow{
		Components: []dgo.MessageComponent{
			&dgo.Button{
				Label:    "게임시작",         // 버튼 텍스트
				Style:    dgo.PrimaryButton,  // 버튼 스타일
				CustomID: "wink_start_button",     // 버튼 클릭 시 처리할 ID
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

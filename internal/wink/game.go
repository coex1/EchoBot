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

// external packages
import (
	dgo "github.com/bwmarrin/discordgo"
)

var (
  // 윙크 받아서 버튼은 클릭 한 사용자들 
	checkedUsers      = make(map[string]bool)
	totalParticipants = 0
	messageIDMap      = make(map[string]string)
)

func FollowUpHandler(s *dgo.Session, i *dgo.InteractionCreate) {
	userID := i.Member.User.ID
	userGlobalName := i.Member.User.GlobalName
	action := ""

	// 버튼 클릭에 따라 상태 업데이트
	switch i.MessageComponentData().CustomID {
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
	checkedCount := CountCheckedUsers()

	// 기존 메시지 업데이트 (Followup 메시지 수정)
	messageID := messageIDMap[i.GuildID]
	if messageID == "" {
		log.Println("No message found to update")
		return
	}

	// 체크된 유저 및 체크되지 않은 유저 목록 생성
	var checkedUsersList, uncheckedUsersList string
	for _, id := range selectedUsersMap[i.GuildID] {
		// 유저 정보를 가져오기
		member, err := s.GuildMember(i.GuildID, id)
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
	var embed *discordgo.MessageEmbed
	if checkedCount == totalParticipants-1 {
		lastUserName := uncheckedUsersList
		embed = &discordgo.MessageEmbed{
			Title: "마지막 남은 사람!\n",
			Description: fmt.Sprintf(
				"%s님, 당신이 마지막 사람입니다.\n\n왕일 것 같은 사람을 지목해주세요!", strings.ReplaceAll(lastUserName, "\n", ""),
			),
			Color: 0xff0000, // 다른 색으로 표시
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title: "게임 진행 중...\n",
			Description: fmt.Sprintf(
				"윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** %d / %d\n\n**확인한 유저 :**\n%s\n**확인하지 못한 유저 :**\n%s",
				checkedCount, totalParticipants, checkedUsersList, uncheckedUsersList,
			),
			Color: 0x00ff00,
		}
	}

	// 상호작용 응답 지연 후 아래 메시지 수정 진행
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	// 메시지 수정
	content := fmt.Sprintf("'%s'이(가) %s했습니다.\n", userGlobalName, action)
	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    i.ChannelID,
		ID:         messageID,
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Content:    &content,
		Components: &i.Message.Components, // 기존 버튼 컴포넌트 유지
	})
	if err != nil {
		log.Println("Error updating message:", err)
		return
	}
}

// 버튼 포함 임베드 메시지 생성
func CreateFollowUpMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	totalParticipants = len(selectedUsersMap[i.GuildID])

	embed := &discordgo.MessageEmbed{
		Title:       "게임 시작!",
		Description: "윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** 0 / " + strconv.Itoa(totalParticipants),
		Color:       0x00ff00,
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "V",
					Style:    discordgo.SuccessButton,
					CustomID: "check",
				},
				&discordgo.Button{
					Label:    "X",
					Style:    discordgo.DangerButton,
					CustomID: "cancel",
				},
			},
		},
	}

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		log.Println("Error sending follow-up message:", err)
	}

	messageIDMap[i.GuildID] = msg.ID
}

func CountCheckedUsers() int {
	count := 0
	for _, checked := range checkedUsers {
		if checked {
			count++
		}
	}
	return count
}

// Select 버튼이 눌렸을 때 선택된 멤버들을 처리하는 핸들러
func winkStartButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		discord.SendDM(s, id, message)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	winkFollowUpMessage(s, i)
}

func winkFollowUpHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	userGlobalName := i.Member.User.GlobalName
	action := ""

	// 버튼 클릭에 따라 상태 업데이트
	switch i.MessageComponentData().CustomID {
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
	checkedCount := countCheckedUsers()

	// 기존 메시지 업데이트 (Followup 메시지 수정)
	messageID := messageIDMap[i.GuildID]
	if messageID == "" {
		log.Println("No message found to update")
		return
	}

	// 체크된 유저 및 체크되지 않은 유저 목록 생성
	var checkedUsersList, uncheckedUsersList string
	for _, id := range selectedUsersMap[i.GuildID] {
		// 유저 정보를 가져오기
		member, err := s.GuildMember(i.GuildID, id)
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
	var embed *discordgo.MessageEmbed
	if checkedCount == totalParticipants-1 {
		lastUserName := uncheckedUsersList
		embed = &discordgo.MessageEmbed{
			Title: "마지막 남은 사람!\n",
			Description: fmt.Sprintf(
				"%s님, 당신이 마지막 사람입니다.\n\n왕일 것 같은 사람을 지목해주세요!", strings.ReplaceAll(lastUserName, "\n", ""),
			),
			Color: 0xff0000, // 다른 색으로 표시
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title: "게임 진행 중...\n",
			Description: fmt.Sprintf(
				"윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** %d / %d\n\n**확인한 유저 :**\n%s\n**확인하지 못한 유저 :**\n%s",
				checkedCount, totalParticipants, checkedUsersList, uncheckedUsersList,
			),
			Color: 0x00ff00,
		}
	}

	// 상호작용 응답 지연 후 아래 메시지 수정 진행
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	// 메시지 수정
	content := fmt.Sprintf("'%s'이(가) %s했습니다.\n", userGlobalName, action)
	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    i.ChannelID,
		ID:         messageID,
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Content:    &content,
		Components: &i.Message.Components, // 기존 버튼 컴포넌트 유지
	})
	if err != nil {
		log.Println("Error updating message:", err)
		return
	}
}

// 버튼 포함 임베드 메시지 생성
func winkFollowUpMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	totalParticipants = len(selectedUsersMap[i.GuildID])

	embed := &discordgo.MessageEmbed{
		Title:       "게임 시작!",
		Description: "윙크를 받으셨으면 V 버튼을 클릭 해 주세요!\n\n실수로 V 했을 경우 X 버튼으로 취소 해 주세요!\n\n**현재 윙크 받은 사람 수 :** 0 / " + strconv.Itoa(totalParticipants),
		Color:       0x00ff00,
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "V",
					Style:    discordgo.SuccessButton,
					CustomID: "check",
				},
				&discordgo.Button{
					Label:    "X",
					Style:    discordgo.DangerButton,
					CustomID: "cancel",
				},
			},
		},
	}

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		log.Println("Error sending follow-up message:", err)
	}

	messageIDMap[i.GuildID] = msg.ID
}

func countCheckedUsers() int {
	count := 0
	for _, checked := range checkedUsers {
		if checked {
			count++
		}
	}
	return count
}

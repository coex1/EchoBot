package discord

// system packages
import (
	"log"
)

// internal imports
import (
	"github.com/coex1/EchoBot/internal/wink"
	"github.com/coex1/EchoBot/internal/mafia"
)

// external packages
import (
	dgo "github.com/bwmarrin/discordgo"
)

const MAX_MEMBER_GET int = 50
const QUERY_STRING string = ""

func handleApplicationCommand(s *dgo.Session, event *dgo.InteractionCreate) {
  switch event.ApplicationCommandData().Name {
  case "wink":
    winkCommandHandle(s, event)
  case "mafia":
    mafiaCommandHandle(s, event)
  }
}

func winkCommandHandle(s *dgo.Session, event *dgo.InteractionCreate) {
	members, err := s.GuildMembers(event.GuildID, QUERY_STRING, MAX_MEMBER_GET)

  // if getting members failed
  if err != nil {
    log.Fatalf("Failed getting members [%v]", err)
    return
  }

	// User 목록으로 SelectMenu 생성
	var options []dgo.SelectMenuOption
	for _, member := range members {
		// member.User.ID와 member.User.Username을 사용하여 옵션 생성
		if !member.User.Bot {
			options = append(options, dgo.SelectMenuOption{
				Label: member.User.GlobalName,
				Value: member.User.ID,
			})
		}
	}
	if game == "wink" {
		MinValues = 3
	} else if game == "mafia" {
		MinValues = 5
	}
	MaxValues = len(options)

	// SelectMenu와 ActionRow 설정
	selectMenu := discordgo.SelectMenu{
		CustomID:    "user_select_menu",
		Placeholder: "Select a user...",
		MinValues:   &MinValues,
		MaxValues:   MaxValues,
		Options:     options,
	}
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	}

	// start_button
	buttonRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.Button{
				Label:    "시작",                    // 버튼 텍스트
				Style:    discordgo.PrimaryButton, // 버튼 스타일
				CustomID: "start_button",          // 버튼 클릭 시 처리할 ID
			},
		},
	}

	// 드롭다운 메뉴와 버튼을 포함한 메시지 전송
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				actionRow,
				buttonRow,
			},
		},
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	selectedUsersMap[i.GuildID] = make([]string, 0)
}

func mafiaCommandHandle(s *discordgo.Session, i *discordgo.InteractionCreate, game string) {
	// 길드 멤버 가져오기
	members, err := s.GuildMembers(i.GuildID, "", 25)
	if err != nil {
		log.Println("Error fetching members:", err)
		return
	}

	// User 목록으로 SelectMenu 생성
	var options []discordgo.SelectMenuOption
	for _, member := range members {
		// member.User.ID와 member.User.Username을 사용하여 옵션 생성
		if !member.User.Bot {
			options = append(options, discordgo.SelectMenuOption{
				Label: member.User.GlobalName,
				Value: member.User.ID,
			})
		}
	}
	if game == "wink" {
		MinValues = 3
	} else if game == "mafia" {
		MinValues = 5
	}
	MaxValues = len(options)

	// SelectMenu와 ActionRow 설정
	selectMenu := discordgo.SelectMenu{
		CustomID:    "user_select_menu",
		Placeholder: "Select a user...",
		MinValues:   &MinValues,
		MaxValues:   MaxValues,
		Options:     options,
	}
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	}

	// start_button
	buttonRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.Button{
				Label:    "시작",                    // 버튼 텍스트
				Style:    discordgo.PrimaryButton, // 버튼 스타일
				CustomID: "start_button",          // 버튼 클릭 시 처리할 ID
			},
		},
	}

	// 드롭다운 메뉴와 버튼을 포함한 메시지 전송
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				actionRow,
				buttonRow,
			},
		},
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	selectedUsersMap[i.GuildID] = make([]string, 0)
}

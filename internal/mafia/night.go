package mafia

import (
	"log"
	"strconv"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
)

func Night_Message(s *dgo.Session, guild *data.Guild) {
	AliveUsersID := Vote_Reset(guild)
	day := strconv.Itoa(guild.Mafia.Day)

	for _, player := range guild.Mafia.Players {
		if player.IsAlive {
			switch player.Role {
			case "Mafia":
				message := &dgo.MessageSend{
					Embeds: []*dgo.MessageEmbed{
						{
							Title:       day + "일 차 밤입니다.",
							Description: "제거할 대상을 선택한 후 '제거하기' 버튼을 눌러주세요!",
							Color:       0xC87C00,
						},
					},
					Components: []dgo.MessageComponent{
						dgo.ActionsRow{
							Components: []dgo.MessageComponent{
								&dgo.SelectMenu{
									MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
									CustomID:    "mafia_Mafia_listUpdate",
									Placeholder: "한 명을 선택하세요",
									MaxValues:   1,
									Options:     AliveUsersID,
								},
							},
						},
						dgo.ActionsRow{
							Components: []dgo.MessageComponent{
								&dgo.Button{
									Label:    "제거하기",
									Style:    dgo.PrimaryButton,
									CustomID: "mafia_Mafia_Button",
								},
							},
						},
					},
				}
				_, err := s.ChannelMessageSendComplex(player.DMChannelID, message)
				if err != nil {
					log.Printf("Failed to send confirmation DM to user %s: %v\n", player.ID, err)
				}
			case "Police":
				message := &dgo.MessageSend{
					Embeds: []*dgo.MessageEmbed{
						{
							Title:       day + "일 차 밤입니다.",
							Description: "직업을 확인하고 싶은 대상을 선택한 후 '조사하기' 버튼을 눌러주세요!",
							Color:       0xC87C00,
						},
					},
					Components: []dgo.MessageComponent{
						dgo.ActionsRow{
							Components: []dgo.MessageComponent{
								&dgo.SelectMenu{
									MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
									CustomID:    "mafia_Police_listUpdate",
									Placeholder: "한 명을 선택하세요",
									MaxValues:   1,
									Options:     AliveUsersID,
								},
							},
						},
						dgo.ActionsRow{
							Components: []dgo.MessageComponent{
								&dgo.Button{
									Label:    "조사하기",
									Style:    dgo.PrimaryButton,
									CustomID: "mafia_Police_Button",
								},
							},
						},
					},
				}
				_, err := s.ChannelMessageSendComplex(player.DMChannelID, message)
				if err != nil {
					log.Printf("Failed to send confirmation DM to user %s: %v\n", player.ID, err)
				}
			case "Doctor":
				message := &dgo.MessageSend{
					Embeds: []*dgo.MessageEmbed{
						{
							Title:       day + "일 차 밤입니다.",
							Description: "살리고 싶은 대상을 선택한 후 '치료하기' 버튼을 눌러주세요!",
							Color:       0xC87C00,
						},
					},
					Components: []dgo.MessageComponent{
						dgo.ActionsRow{
							Components: []dgo.MessageComponent{
								&dgo.SelectMenu{
									MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
									CustomID:    "mafia_Doctor_listUpdate",
									Placeholder: "한 명을 선택하세요",
									MaxValues:   1,
									Options:     AliveUsersID,
								},
							},
						},
						dgo.ActionsRow{
							Components: []dgo.MessageComponent{
								&dgo.Button{
									Label:    "치료하기",
									Style:    dgo.PrimaryButton,
									CustomID: "mafia_Doctor_Button",
								},
							},
						},
					},
				}
				_, err := s.ChannelMessageSendComplex(player.DMChannelID, message)
				if err != nil {
					log.Printf("Failed to send confirmation DM to user %s: %v\n", player.ID, err)
				}
			case "Citizen":
				message := &dgo.MessageSend{
					Embeds: []*dgo.MessageEmbed{
						{
							Title:       day + "일 차 밤입니다.",
							Description: "졸려",
							Color:       0xC87C00,
						},
					},
					Components: []dgo.MessageComponent{
						dgo.ActionsRow{
							Components: []dgo.MessageComponent{
								&dgo.SelectMenu{
									MenuType:    dgo.SelectMenuType(dgo.SelectMenuComponent),
									CustomID:    "mafia_Doctor_listUpdate",
									Placeholder: "한 명을 선택하세요",
									MaxValues:   1,
									Options:     AliveUsersID,
								},
							},
						},
						dgo.ActionsRow{
							Components: []dgo.MessageComponent{
								&dgo.Button{
									Label:    "치료하기",
									Style:    dgo.PrimaryButton,
									CustomID: "mafia_Doctor_Button",
								},
							},
						},
					},
				}
				_, err := s.ChannelMessageSendComplex(player.DMChannelID, message)
				if err != nil {
					log.Printf("Failed to send confirmation DM to user %s: %v\n", player.ID, err)
				}
			}
		}
	}
}

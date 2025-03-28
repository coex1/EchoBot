package mafia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"
)

func announceNightResult(s *dgo.Session, guild *data.Guild) {
	players := guild.Mafia.Players

	var result string

	// Save or Not
	if guild.Mafia.MafiaTarget == "" || players[guild.MafiaTarget] == nil {
		log.Printf("nil")
	} else {
		if guild.Mafia.MafiaTarget == guild.Mafia.DoctorTarget {
			result = fmt.Sprintf("**%s** 님이 공격당했지만, 의사의 치료로 생존했습니다!", players[guild.Mafia.MafiaTarget].GlobalName)
		} else {
			result = fmt.Sprintf("**%s** 님이 마피아의 공격으로 사망했습니다.", players[guild.Mafia.MafiaTarget].GlobalName)
			players[guild.Mafia.MafiaTarget].IsAlive = false
		}

		nightResult := &dgo.MessageEmbed{
			Title:       "밤이 지나고...",
			Description: result,
			Color:       0x992D22,
		}
		for id := range players {
			general.SendComplexDM(s, id, &dgo.MessageSend{
				Embeds: []*dgo.MessageEmbed{nightResult},
			})
		}
	}
}

func isNightActionAllDone(guild *data.Guild) bool {
	// 생존해 있는 마피아, 경찰, 의사가 대상 선택했는지 확인
	for _, p := range guild.Mafia.Players {
		if p.Role == "Police" && !p.IsAlive {
			guild.Mafia.NightActionDone["Police"] = true
		}
		if p.Role == "Doctor" && !p.IsAlive {
			guild.Mafia.NightActionDone["Doctor"] = true
		}
	}

	if !guild.Mafia.NightActionDone["Mafia"] {
		return false
	}
	if !guild.Mafia.NightActionDone["Police"] {
		return false
	}
	if !guild.Mafia.NightActionDone["Doctor"] {
		return false
	}
	return true
}

func Night_Message(s *dgo.Session, guild *data.Guild) {
	AliveUsersID := Reset(guild)
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
				general.SendComplexDM(s, player.ID, message)
				// _, err := s.ChannelMessageSendComplex(player.DMChannelID, message)
				// if err != nil {
				// 	log.Printf("Failed to send confirmation DM to user %s: %v\n", player.ID, err)
				// }
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
				general.SendComplexDM(s, player.ID, message)
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
				general.SendComplexDM(s, player.ID, message)
			case "Citizen":
				message := &dgo.MessageSend{
					Embeds: []*dgo.MessageEmbed{
						{
							Title:       day + "일 차 밤입니다.",
							Description: "아래 문장을 입력하세요",
							Color:       0xC87C00,
						},
					},
				}
				general.SendComplexDM(s, player.ID, message)
			}
		}
	}
}

func Mafia_listUpdate(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.MafiaTargetMap[i.User.ID] = i.MessageComponentData().Values[0]
}
func Mafia_Skill_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	players := guild.Mafia.Players

	mafiaAlive := 0
	var aliveMafiaList []string
	for _, id := range players {
		if id.IsAlive && id.Role == "Mafia" {
			mafiaAlive += 1
			aliveMafiaList = append(aliveMafiaList, id.ID)
		}
	}

	if len(guild.Mafia.MafiaTargetMap) == mafiaAlive {
		// 투표 집계 업데이트
		voteCount := make(map[string]int)
		for _, targetID := range guild.Mafia.MafiaTargetMap {
			voteCount[targetID] += 1
		}

		// 최다 득표 집계
		maxVotes := 0
		var topVoted []string
		for id, count := range voteCount {
			if count > maxVotes {
				maxVotes = count
				topVoted = []string{id}
			} else if count == maxVotes {
				topVoted = append(topVoted, id)
			}
		}
		guild.Mafia.MafiaTarget = topVoted[0]
		if len(topVoted) > 1 {
			guild.Mafia.MafiaTarget = aliveMafiaList[general.Random(0, len(topVoted)-1)]
		}

		guild.Mafia.NightActionDone["Mafia"] = true

		if isNightActionAllDone(guild) {
			time.Sleep(3 * time.Second)
			announceNightResult(s, guild)

			if isGameOver(guild) {
				gameEndingMessage(s, guild)
			} else {
				time.Sleep(3 * time.Second)
				Day_Message(s, i, guild)
			}
		}
	}
}

func Police_listUpdate(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.PoliceTarget = i.MessageComponentData().Values[0]
}
func Police_Skill_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	players := guild.Mafia.Players
	var message string
	if guild.Mafia.Players[guild.Mafia.PoliceTarget].Role == "Mafia" {
		message = "맞습니다!"
	} else {
		message = "아닙니다!"
	}

	policeSkill := &dgo.MessageEmbed{
		Title:       "조사 결과",
		Description: fmt.Sprintf("**%s** 님은 마피아가 ", players[guild.Mafia.PoliceTarget].GlobalName) + message,
		Color:       0x3498db,
	}
	general.SendComplexDM(s, i.User.ID, &dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			policeSkill,
		},
	})

	guild.Mafia.NightActionDone["Police"] = true

	if isNightActionAllDone(guild) {
		time.Sleep(3 * time.Second)
		announceNightResult(s, guild)

		if isGameOver(guild) {
			gameEndingMessage(s, guild)
		} else {
			time.Sleep(3 * time.Second)
			Day_Message(s, i, guild)
		}
	}
}

func Doctor_listUpdate(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.DoctorTarget = i.MessageComponentData().Values[0]
}
func Doctor_Skill_Button(s *dgo.Session, i *dgo.InteractionCreate, guild *data.Guild) {
	guild.Mafia.NightActionDone["Doctor"] = true

	if isNightActionAllDone(guild) {
		time.Sleep(3 * time.Second)
		announceNightResult(s, guild)

		if isGameOver(guild) {
			gameEndingMessage(s, guild)
		} else {
			time.Sleep(3 * time.Second)
			Day_Message(s, i, guild)
		}
	}
}

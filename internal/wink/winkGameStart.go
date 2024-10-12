package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

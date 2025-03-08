package mafia

import (

	// system

	"math/rand"
	"time"

	// external package
	dgo "github.com/bwmarrin/discordgo"

	// internal package

	"github.com/coex1/EchoBot/internal/data"
	"github.com/coex1/EchoBot/internal/general"
)

// 플레이어 역할 배정
func Role_Message(s *dgo.Session, guild *data.Guild) {

	// embed for Mafia
	embedMafia := dgo.MessageEmbed{
		Title:       "당신은 **마피아**입니다!",
		Description: "밤마다 시민을 처치할 수 있습니다.",
		Color:       0xFFD800,
	}
	// data for Mafia
	dataMafia := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedMafia,
		},
	}

	// embed for Police
	embedPolice := dgo.MessageEmbed{
		Title:       "당신은 **경찰**입니다!",
		Description: "밤마다 한 명의 신원을 확인할 수 있습니다.",
		Color:       0xC87C00,
	}
	// data for Police
	dataPolice := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedPolice,
		},
	}

	// embed for Doctor
	embedDoctor := dgo.MessageEmbed{
		Title:       "당신은 **의사**입니다!",
		Description: "밤마다 마피아로부터 한 명을 보호할 수 있습니다.",
		Color:       0xC87C00,
	}
	// data for Doctor
	dataDoctor := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedDoctor,
		},
	}

	//embed for Citizen
	embedCitizen := dgo.MessageEmbed{
		Title:       "당신은 **시민**입니다!",
		Description: "마피아를 찾아내세요.",
		Color:       0xC87C00,
	}
	// data for Citizen
	dataCitizen := dgo.MessageSend{
		Embeds: []*dgo.MessageEmbed{
			&embedCitizen,
		},
	}

	// shuffle algorithm
	shuffled := make([]string, len(guild.Mafia.SelectedUsersID))
	copy(shuffled, guild.Mafia.SelectedUsersID)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	mafiaIDs := shuffled[:guild.Mafia.NumMafia]
	shuffled = shuffled[guild.Mafia.NumMafia:]
	policeIDs := shuffled[:guild.Mafia.NumPolice]
	shuffled = shuffled[guild.Mafia.NumPolice:]
	doctorIDs := shuffled[:guild.Mafia.NumDoctor]

	for _, player := range guild.Mafia.Players {
		switch {
		case general.Contains(mafiaIDs, player.ID):
			dmChannelID, _ := general.Mafia_SendComplexDM(s, player.ID, &dataMafia)
			player.DMChannelID = dmChannelID
			player.Role = "Mafia"

		case general.Contains(policeIDs, player.ID):
			dmChannelID, _ := general.Mafia_SendComplexDM(s, player.ID, &dataPolice)
			player.DMChannelID = dmChannelID
			player.Role = "Police"

		case general.Contains(doctorIDs, player.ID):
			dmChannelID, _ := general.Mafia_SendComplexDM(s, player.ID, &dataDoctor)
			player.DMChannelID = dmChannelID
			player.Role = "Doctor"

		default:
			dmChannelID, _ := general.Mafia_SendComplexDM(s, player.ID, &dataCitizen)
			player.DMChannelID = dmChannelID
			player.Role = "Citizen"
		}
	}
}
